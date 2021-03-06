package analysis

import "time"

type Report struct {
	Timestamp            time.Time
	DomainNumber         int
	TrackerNumber        int
	CookieData           *ReportCookieData
	ThirdPartyCookieData *ReportSpecificCookieData
	TrackerCookieData    *ReportSpecificCookieData
}

type ReportCookieData struct {
	Amount               int
	AveragePerSite       float32
	MedianPerSite        float32
	WebsitesAboveAverage int
	WebsitesNoUsage      int
}

type ReportSpecificCookieData struct {
	Amount               int
	ShareInCookies       float32
	AveragePerSite       float32
	MedianPerSite        float32
	WebsitesAboveAverage int
	WebsitesNoUsage      int
}
