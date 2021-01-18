# website-cookie-analyzer
This project aims to provide websites statistics about cookie usage.

# Description

This project has been developed as part of a term paper which analyzes Cookie usage on websites.

# Prerequisites

In order to run this application, the **[Chromium Embedded Framework](https://bitbucket.org/chromiumembedded/cef/)** is 
required.

# Download

The latest binary of this software can be found at the releases page of this project.

# Usage

In order to use this program, the user has to provide a list of websites and a list of trackers. This list of trackers
is required to detect possible tracking cookies.

## List of websites

Syntax:
```csv
google.de
github.com
facebook.com
studenten.ba-rm.de
```

Bash command example in order to fetch Alexa Top 100 Sites and write to `websites.csv` with a compliant file format:

```bash
COUNTRY=DE curl -H "x-api-key: <token>" "https://ats.api.alexa.com/api?Action=Topsites&Count=100&CountryCode=${COUNTRY}&ResponseGroup=Country&Start=1&Output=json" | jq -r '.Ats.Results.Result.Alexa.TopSites.Country.Sites.Site[].DataUrl' > websites.csv
```

## List of trackers

Syntax:

```
# this is a comment
0.0.0.0 this-is-a-tracker.com
0.0.0.0 no-good-website.com
0.0.0.0 malicious-website.com
127.0.0.1 this-line-will-be-ignored
```

Example: https://github.com/StevenBlack/hosts/

# Compiling

In order to compile this binary, the source code has to be cloned and Golang v1.15.* must be installed locally. 
Compiling is done by running the following command:

```bash
go build ./cmd/websitecookieanalyzer/main.go
```

# Algorithm

The program analyzes cookies based on the Domain-Attribute used - in concrete terms, this means that 3rd party cookies 
are detected when the Domain-Attribute does not end with the website`s domain (in order to detect subdomains). If so,
the cookie is marked as a 3rd party one. If the Domain-Attribute is a website listed on the tracker list, the program 
assumes that the cookie is a tracker cookie.

When generating the report, all data is aggregated so no individual website data is included. 
