package fetch

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

type CookieFetchResult struct {
	Date                             time.Time
	ChromeVersion, OSType, UserAgent string
	Cookies                          []*WebsiteCookies
}

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

func FetchCookies(websites []string) *CookieFetchResult {
	result := &CookieFetchResult{
		Date: time.Now(),
	}
	logrus.Infoln("Trying to retrieve Chrome version...")
	fetchChromeVersion(result)
	logrus.WithField("chromeVersion", result.ChromeVersion).Infoln("Retrieved chrome version.")
	logrus.Infoln("Fetching cookies for websites...")
	result.Cookies = make([]*WebsiteCookies, 0)
	for index, website := range websites {
		logrus.WithField("index", index+1).WithField("website", website).Debugln("Scanning cookies of website...")
		websiteCookies, err := fetchCookiesSingleWebsite(website)
		if err != nil {
			logrus.WithError(err).WithField("website", website).Errorln("Could not retrieve cookies for website!")
		} else {
			result.Cookies = append(result.Cookies, &WebsiteCookies{
				WebsiteName: website,
				Cookies:     websiteCookies,
			})
			logrus.WithField("website", website).WithField("cookieAmount", len(websiteCookies)).Debugln("Scanned website cookies!")
		}
	}
	logrus.WithField("cookieFetchWebsites", len(result.Cookies)).Infoln("Fetched cookies for websites.")
	return result
}

func fetchChromeVersion(result *CookieFetchResult) {
	dir, ctx, cancel, err := setupChromeAndNavigate("chrome://version/")
	defer func() {
		if dir == "" {
			return
		}
		time.Sleep(cleanUpDelay)
		if err := os.RemoveAll(dir); err != nil {
			logrus.WithError(err).WithField("chromeDir", dir).Warnln("Could not delete chrome working directory!")
		}
	}()
	defer func() {
		if cancel != nil {
			cancel()
		}
	}()
	if err != nil {
		logrus.WithError(err).Fatalln("Could not setup chrome and navigate to chrome about website!")
		return
	}
	if err = chromedp.Run(ctx, chromedp.Text("#version", &result.ChromeVersion, chromedp.NodeVisible, chromedp.ByID)); err != nil {
		logrus.WithError(err).Fatalln("Could not retrieve chrome version from about site.")
	}
	if err = chromedp.Run(ctx, chromedp.Text("#os_type", &result.OSType, chromedp.NodeVisible, chromedp.ByID)); err != nil {
		logrus.WithError(err).Fatalln("Could not retrieve os from about site.")
	}
	if err = chromedp.Run(ctx, chromedp.Text("#useragent", &result.UserAgent, chromedp.NodeVisible, chromedp.ByID)); err != nil {
		logrus.WithError(err).Fatalln("Could not retrieve user agent from about site.")
	}
	cancel()
	return
}

func setupChromeAndNavigate(url string) (string, context.Context, context.CancelFunc, error) {
	dir, err := setupChromeDirectory()
	if err != nil {
		logrus.WithError(err).Fatalln("Could not setup chrome directory!")
	}
	options := append(chromeExecOptions, chromedp.UserDataDir(dir))
	execCtx, cancel := chromedp.NewExecAllocator(context.Background(), options...)
	ctx, _ := chromedp.NewContext(execCtx)
	if err := chromedp.Run(ctx); err != nil {
		logrus.WithError(err).WithField("url", url).Fatalln("Could not run browser process!")
	}
	timeOutCtx, _ := context.WithTimeout(ctx, navigateTimeout)
	err = chromedp.Run(timeOutCtx, chromedp.Navigate(url))
	if err != nil {
		if err == context.DeadlineExceeded {
			logrus.WithField("url", url).WithField("timeout", navigateTimeout.String()).Warnln("Website load timeout exceeded.")
			err = nil
		} else {
			logrus.WithError(err).WithField("url", url).Errorln("Could not navigate to url!")
			return dir, ctx, cancel, err
		}
	}
	return dir, ctx, cancel, err
}

func fetchCookiesSingleWebsite(website string) ([]*network.Cookie, error) {
	cookies, err := fetchBrowserCookies(website)
	return cookies, err
}

func fetchBrowserCookies(website string) ([]*network.Cookie, error) {
	url := fmt.Sprintf("http://%s", website)
	dir, ctx, cancel, err := setupChromeAndNavigate(url)
	defer func() {
		if dir == "" {
			return
		}
		time.Sleep(cleanUpDelay)
		if err := os.RemoveAll(dir); err != nil {
			logrus.WithError(err).WithField("chromeDir", dir).Warnln("Could not delete chrome working directory!")
		}
	}()
	defer func() {
		if cancel != nil {
			cancel()
		}
	}()
	if err != nil {
		logrus.WithError(err).Warnln("Could not setup chrome and navigate to cookie related website!")
		return nil, err
	}
	time.Sleep(cookieSetWait)
	var cookies []*network.Cookie
	action := chromedp.ActionFunc(func(ctx context.Context) (err error) {
		cookies, err = network.GetAllCookies().Do(ctx)
		if err != nil {
			return err
		}
		return nil
	})
	if err = chromedp.Run(ctx, action); err != nil {
		logrus.WithError(err).WithField("website", website).Errorln("Could not retrieve cookies!")
		return nil, err
	}
	return cookies, nil
}

func setupChromeDirectory() (name string, err error) {
	dir, err := ioutil.TempDir("", tempDirPattern)
	return dir, err
}
