package analysis

import (
	"bufio"
	"github.com/chromedp/cdproto/network"
	"github.com/mmichaelb/website-cookie-analyzer/internal/pkg/websitecookieanalyzer/fetch"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
)

func analyzeTrackerCookies(cookieData []*fetch.WebsiteCookies, report *Report, trackers []string) (thirdPartyCookieData *ReportSpecificCookieData) {
	trackerCookieCheckFunc := func(domain string, cookie *network.Cookie) bool {
		for _, tracker := range trackers {
			// remove leading dots
			trimmedCookieDomain := strings.Trim(cookie.Domain, ".")
			if strings.EqualFold(trimmedCookieDomain, tracker) {
				return true
			}
		}
		return false
	}
	return analyzeSpecificCookies(cookieData, report, trackerCookieCheckFunc)
}

func LoadTrackers(trackersFilepath string) ([]string, error) {
	file, err := os.Open(trackersFilepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	websites := make([]string, 0)
	scanner := bufio.NewScanner(file)

	for {
		ok := scanner.Scan()
		if !ok {
			if err = scanner.Err(); err != nil {
				logrus.WithError(err).Fatalln("Could not read tracker line!")
			}
			break
		}
		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 || line[0] == '#' {
			continue
		}
		lineSplit := strings.SplitN(line, " ", 2)
		ip := lineSplit[0]
		host := lineSplit[1]
		if ip != "0.0.0.0" && host != "0.0.0.0" {
			continue
		}
		websites = append(websites, host)
	}
	return websites, nil
}
