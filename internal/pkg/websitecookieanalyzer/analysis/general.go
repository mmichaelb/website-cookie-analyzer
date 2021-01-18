package analysis

import (
	"github.com/mmichaelb/website-cookie-analyzer/internal/pkg/websitecookieanalyzer/fetch"
	"github.com/mmichaelb/website-cookie-analyzer/internal/pkg/websitecookieanalyzer/util"
	"sort"
	"time"
)

func Analyze(cookieData []*fetch.WebsiteCookies, trackers []string) (report *Report) {
	report = &Report{
		Timestamp:     time.Now(),
		DomainNumber:  len(cookieData),
		TrackerNumber: len(trackers),
		CookieData:    analyzeCookies(cookieData),
	}
	report.ThirdPartyCookieData = analyzeThirdPartyCookies(cookieData, report)
	report.TrackerCookieData = analyzeTrackerCookies(cookieData, report, trackers)
	return report
}

func analyzeCookies(cookieData []*fetch.WebsiteCookies) (reportCookieData *ReportCookieData) {
	reportCookieData = &ReportCookieData{
		Amount:               -1,
		AveragePerSite:       -1,
		MedianPerSite:        -1,
		WebsitesAboveAverage: -1,
		WebsitesNoUsage:      -1,
	}
	reportCookieData.Amount = 0
	reportCookieData.WebsitesNoUsage = 0
	cookieAmountSlice := make([]int, 0)
	for _, elemCookieData := range cookieData {
		cookieAmount := len(elemCookieData.Cookies)
		cookieAmountSlice = append(cookieAmountSlice, cookieAmount)
		reportCookieData.Amount += cookieAmount
		if cookieAmount == 0 {
			reportCookieData.WebsitesNoUsage++
		}
	}
	// retrieve median
	sort.Ints(cookieAmountSlice)
	reportCookieData.MedianPerSite = util.GetMedian(cookieAmountSlice)
	websiteAmount := len(cookieData)
	reportCookieData.AveragePerSite = float32(reportCookieData.Amount) / float32(websiteAmount)
	for index, amount := range cookieAmountSlice {
		if float32(amount) > reportCookieData.AveragePerSite {
			reportCookieData.WebsitesAboveAverage = len(cookieAmountSlice) - index
			break
		}
	}
	return reportCookieData
}
