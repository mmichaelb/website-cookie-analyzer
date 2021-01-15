package main

import (
	"encoding/xml"
	"flag"
	"github.com/mmichaelb/website-cookie-analyzer/internal/pkg/websitecookieanalyzer"
	"github.com/sirupsen/logrus"
	"os"
)

var (
	levelString           = flag.String("level", "info", "Sets the logging level.")
	websitesInputFilepath = flag.String("websitesFile", "./websites.csv", "Sets the file path to the input websites file.")
	cookiesFilepath       = flag.String("cookiesFile", "./cookies.xml", "Sets the file path to the output cookies file.")
	fetchNewCookies       = flag.Bool("fetch", true, "Determines whether the domain input file should be used in order to fetch the cookies.")

	websites []string
	cookies  []*websitecookieanalyzer.WebsiteCookies
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
	logrus.Exit(0)
}

func loadWebsites() {
	logrus.WithField("websitesFile", *websitesInputFilepath).Infoln("Loading website file...")
	var err error
	websites, err = websitecookieanalyzer.LoadWebsites(*websitesInputFilepath)
	if err != nil {
		logrus.WithError(err).Fatalln("Could not read website file!")
	}
	logrus.WithField("websiteCount", len(websites)).Infoln("Successfully loaded website file.")
}

func fetchCookies() {
	logrus.Infoln("Fetching cookies for websites...")
	cookies = websitecookieanalyzer.FetchCookies(websites)
	logrus.WithField("cookieFetchWebsites", len(cookies)).Infoln("Fetched cookies for websites.")
}

func writeCookies() {
	logrus.WithField("cookiesOutputFile", *cookiesFilepath).Infoln("Writing cookies to cookie output file...")
	file, err := os.Create(*cookiesFilepath)
	if err != nil {
		logrus.WithError(err).Fatalln("Could not create cookie output file!")
	}
	defer file.Close()
	xmlBytes, err := xml.MarshalIndent(cookies, "", "  ")
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
