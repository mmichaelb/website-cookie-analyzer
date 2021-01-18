package analysis

import (
	"github.com/chromedp/cdproto/network"
	"github.com/mmichaelb/website-cookie-analyzer/internal/pkg/websitecookieanalyzer/fetch"
	"strings"
)

func analyzeThirdPartyCookies(cookieData []*fetch.WebsiteCookies, report *Report) (thirdPartyCookieData *ReportSpecificCookieData) {
	thirdPartyCheckFunc := func(domain string, cookie *network.Cookie) bool {
		return !strings.HasSuffix(cookie.Domain, domain)
	}
	return analyzeSpecificCookies(cookieData, report, thirdPartyCheckFunc)
}
