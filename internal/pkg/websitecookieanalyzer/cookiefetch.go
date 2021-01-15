package websitecookieanalyzer

import (
	"context"
	"fmt"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"time"
)

const (
	tempDirPattern  = "wca-chrome"
	navigateTimeout = time.Second * 30
	// sets the duration how long to wait until collecting the cookies
	cookieSetWait = time.Second * 5
	// sets the duration how to wait before deleting the temp directory
	cleanUpDelay = time.Second
)

type WebsiteCookies struct {
	WebsiteName string
	Cookies     []*network.Cookie
}

var (
	chromeExecOptions = []chromedp.ExecAllocatorOption{
		chromedp.NoFirstRun,
		chromedp.NoDefaultBrowserCheck,
		chromedp.DisableGPU,
	}
)

func FetchCookies(websites []string) []*WebsiteCookies {
	absoluteCookieList := make([]*WebsiteCookies, 0)
	for _, website := range websites {
		cookies, err := fetchCookiesSingleWebsite(website)
		if err != nil {
			logrus.WithError(err).WithField("website", website).Errorln("Could not retrieve cookies for website!")
		} else {
			absoluteCookieList = append(absoluteCookieList, cookies)
		}
	}
	return absoluteCookieList
}

func fetchCookiesSingleWebsite(website string) (*WebsiteCookies, error) {
	dir, err := setupChromeDirectory()
	if err != nil {
		logrus.WithError(err).Fatalln("Could not setup chrome directory!")
	}
	defer func() {
		if err := os.RemoveAll(dir); err != nil {
			logrus.WithError(err).WithField("chromeDir", dir).Warnln("Could not delete chrome working directory!")
		}
	}()
	cookies, err := fetchBrowserCookies(website, dir)
	time.Sleep(cleanUpDelay)
	return cookies, err
}

func fetchBrowserCookies(website string, dir string) (*WebsiteCookies, error) {
	options := append(chromeExecOptions, chromedp.UserDataDir(dir))
	execCtx, cancel := chromedp.NewExecAllocator(context.Background(), options...)
	defer cancel()
	ctx, cancel := chromedp.NewContext(execCtx)
	defer cancel()
	if err := chromedp.Run(ctx); err != nil {
		logrus.WithError(err).WithField("website", website).Fatalln("Could not run browser process!")
	}
	url := fmt.Sprintf("http://%s", website)
	timeOutCtx, cancel := context.WithTimeout(ctx, navigateTimeout)
	defer cancel()
	err := chromedp.Run(timeOutCtx, chromedp.Navigate(url))
	if err != nil {
		if err == context.DeadlineExceeded {
			logrus.WithField("website", website).WithField("timeout", navigateTimeout.String()).Warnln("Website load timeout exceeded.")
		} else {
			logrus.WithError(err).WithField("website", website).Errorln("Could not navigate to website!")
			return nil, err
		}
	}
	time.Sleep(cookieSetWait)
	result := &WebsiteCookies{
		WebsiteName: website,
	}
	action := chromedp.ActionFunc(func(ctx context.Context) error {
		cookies, err := network.GetAllCookies().Do(ctx)
		if err != nil {
			return err
		}
		result.Cookies = cookies
		return nil
	})
	if err = chromedp.Run(ctx, action); err != nil {
		logrus.WithError(err).WithField("website", website).Errorln("Could not retrieve cookies!")
		return nil, err
	}
	return result, nil
}

func setupChromeDirectory() (name string, err error) {
	dir, err := ioutil.TempDir("", tempDirPattern)
	return dir, err
}
