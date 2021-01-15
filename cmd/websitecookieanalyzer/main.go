package main

import (
	"flag"
	"github.com/mmichaelb/website-cookie-analyzer/internal/pkg/websitecookieanalyzer"
	"github.com/sirupsen/logrus"
)

var (
	levelString           = flag.String("level", "info", "Sets the logging level.")
	websitesInputFilepath = flag.String("websitesFile", "./websites.csv", "Sets the file path to the input websites file.")

	websites []string
)

func main() {
	flag.Parse()
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceQuote:       true,
		QuoteEmptyFields: true,
	})
	if level, err := logrus.ParseLevel(*levelString); err != nil {
		logrus.WithError(err).WithField("levelString", *levelString).Warnln("Could not parse custom level string. Falling back to default log level!")
	} else {
		logrus.SetLevel(level)
	}
	loadWebsites()
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
