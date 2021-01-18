package analysis

import (
	"github.com/chromedp/cdproto/network"
	"github.com/mmichaelb/website-cookie-analyzer/internal/pkg/websitecookieanalyzer/fetch"
	"github.com/mmichaelb/website-cookie-analyzer/internal/pkg/websitecookieanalyzer/util"
	"sort"
)

func analyzeSpecificCookies(cookieData []*fetch.WebsiteCookies, report *Report, checkFunc func(domain string, cookie *network.Cookie) bool) (specificCookieData *ReportSpecificCookieData) {
	specificCookieData = &ReportSpecificCookieData{
		Amount:               -1,
		ShareInCookies:       -1,
		AveragePerSite:       -1,
		MedianPerSite:        -1,
		WebsitesAboveAverage: -1,
		WebsitesNoUsage:      -1,
	}
	specificCookieData.Amount = 0
	specificCookieData.WebsitesNoUsage = 0
	specificCookieAmountSlice := make([]int, 0)
	for _, elemCookieData := range cookieData {
		domain := elemCookieData.WebsiteName
		specificCookieCount := 0
		for _, cookie := range elemCookieData.Cookies {
			// check if it is a specific cookie
			specificCookie := checkFunc(domain, cookie)
			if specificCookie {
				// increase local specific cookie count
				specificCookieCount++
				// increase total specific cookie count
				specificCookieData.Amount++
			}
		}
		specificCookieAmountSlice = append(specificCookieAmountSlice, specificCookieCount)
		if specificCookieCount == 0 {
			specificCookieData.WebsitesNoUsage++
		}
	}
	specificCookieData.ShareInCookies = float32(specificCookieData.Amount) / float32(report.CookieData.Amount)
	websiteAmount := len(cookieData)
	specificCookieData.AveragePerSite = float32(specificCookieData.Amount) / float32(websiteAmount)
	sort.Ints(specificCookieAmountSlice)
	specificCookieData.MedianPerSite = util.GetMedian(specificCookieAmountSlice)
	for index, amount := range specificCookieAmountSlice {
		if float32(amount) > specificCookieData.AveragePerSite {
			specificCookieData.WebsitesAboveAverage = len(specificCookieAmountSlice) - index
			break
		}
	}
	if specificCookieData.WebsitesAboveAverage == -1 {
		panic("there should be at least one cookie data elem which is above average")
	}
	return specificCookieData
}
