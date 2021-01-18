package main

import (
	"encoding/xml"
	"flag"
	"github.com/mmichaelb/website-cookie-analyzer/internal/pkg/websitecookieanalyzer/analysis"
	"github.com/mmichaelb/website-cookie-analyzer/internal/pkg/websitecookieanalyzer/fetch"
	"github.com/sirupsen/logrus"
	"os"
)

var (
	levelString           = flag.String("level", "info", "Sets the logging level.")
	websitesInputFilepath = flag.String("websitesFile", "./websites.csv", "Sets the file path to the input websites file.")
	cookiesFilepath       = flag.String("cookiesFile", "./cookies.xml", "Sets the file path to the output cookies file.")
	fetchNewCookies       = flag.Bool("fetch", true, "Determines whether the domain input file should be used in order to fetch the cookies.")
	trackersInputFilepath = flag.String("trackersFile", "./trackers.csv", "Sets the file path to the input trackers file.")
	reportOutputFilepath  = flag.String("reportFile", "./report.xml", "Sets the file path to the output xml file.")

	websites          []string
	cookieFetchResult *fetch.CookieFetchResult
	trackers          []string
)

func main() {
	flag.Parse()
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceQuote:       true,
		QuoteEmptyFields: true,
	})
	logrus.RegisterExitHandler(func() {
		logrus.Infoln("Application stopping. Goodbye!")
	})
	if level, err := logrus.ParseLevel(*levelString); err != nil {
		logrus.WithError(err).WithField("levelString", *levelString).Warnln("Could not parse custom level string. Falling back to default log level!")
	} else {
		logrus.SetLevel(level)
	}
	if *fetchNewCookies {
		loadWebsites()
		fetchCookies()
		writeCookies()
	}
	readCookies()
	loadTrackers()
	runAndSaveAnalysis()
	logrus.Exit(0)
}

func loadWebsites() {
	logrus.WithField("websitesFile", *websitesInputFilepath).Infoln("Loading website file...")
	var err error
	websites, err = fetch.LoadWebsites(*websitesInputFilepath)
	if err != nil {
		logrus.WithError(err).Fatalln("Could not read website file!")
	}
	logrus.WithField("websiteCount", len(websites)).Infoln("Successfully loaded website file.")
}

func fetchCookies() {
	cookieFetchResult = fetch.FetchCookies(websites)
}

func writeCookies() {
	logrus.WithField("cookiesFile", *cookiesFilepath).Infoln("Writing cookies to cookie output file...")
	file, err := os.Create(*cookiesFilepath)
	if err != nil {
		logrus.WithError(err).Fatalln("Could not create cookie output file!")
	}
	defer file.Close()
	xmlBytes, err := xml.MarshalIndent(cookieFetchResult, "", "  ")
	if err != nil {
		logrus.WithError(err).WithField("cookiesOutputFile", *cookiesFilepath).Fatalln("Could not encode cookie output!")
	}
	if _, err = file.WriteString(xml.Header); err != nil {
		logrus.WithError(err).WithField("cookiesOutputFile", *cookiesFilepath).Fatalln("Could not write XML Header!")
	}
	if _, err = file.Write(xmlBytes); err != nil {
		logrus.WithError(err).WithField("cookiesOutputFile", *cookiesFilepath).Fatalln("Could not write XML body!")
	}
	logrus.Infoln("Successfully written cookies for websites.")
}

func readCookies() {
	logrus.WithField("cookiesFile", *cookiesFilepath).Infoln("Reading cookies from cookie file...")
	file, err := os.Open(*cookiesFilepath)
	if err != nil {
		logrus.WithError(err).Fatalln("Could not open cookie file!")
	}
	defer file.Close()
	decoder := xml.NewDecoder(file)
	cookieFetchResult = &fetch.CookieFetchResult{}
	if err = decoder.Decode(cookieFetchResult); err != nil {
		logrus.WithError(err).Fatalln("Could not load cookies from cookie file!")
	}
	test := cookieFetchResult
	logrus.WithField("websiteAmount", len(test.Cookies)).Infoln("Read cookies from cookie file.")
}

func loadTrackers() {
	logrus.WithField("trackersFile", *trackersInputFilepath).Infoln("Loading trackers file...")
	var err error
	trackers, err = analysis.LoadTrackers(*trackersInputFilepath)
	if err != nil {
		logrus.WithError(err).Fatalln("Could not read trackers file!")
	}
	logrus.WithField("trackerCount", len(trackers)).Infoln("Successfully loaded trackers file.")
}

func runAndSaveAnalysis() {
	logrus.Infoln("Running cookie analysis...")
	report := analysis.Analyze(cookieFetchResult.Cookies, trackers)
	logrus.Infoln("Successfully ran cookie analysis!")
	logrus.WithField("reportOutputFile", *reportOutputFilepath).Infoln("Writing cookie analysis report to report file...")
	file, err := os.Create(*reportOutputFilepath)
	if err != nil {
		logrus.WithError(err).Fatalln("Could not create cookie analysis file!")
	}
	defer file.Close()
	xmlBytes, err := xml.MarshalIndent(report, "", "  ")
	if err != nil {
		logrus.WithError(err).WithField("reportOutputFile", *reportOutputFilepath).Fatalln("Could not encode cookie analysis output!")
	}
	if _, err = file.WriteString(xml.Header); err != nil {
		logrus.WithError(err).WithField("reportOutputFile", *reportOutputFilepath).Fatalln("Could not write XML Header!")
	}
	if _, err = file.Write(xmlBytes); err != nil {
		logrus.WithError(err).WithField("reportOutputFile", *reportOutputFilepath).Fatalln("Could not write XML body!")
	}
	logrus.Infoln("Successfully written cookie analysis report to file.")
}
