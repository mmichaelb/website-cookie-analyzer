package analysis

import (
	"bufio"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
)

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
