package main

import (
	"os"
	"regexp"
	"strings"
)

var (
	config        ATSExporterConfig
	defaultConfig = ATSExporterConfig{
		ATSURL: "http://localhost:80/_stats",
	}
)

type ATSExporterConfig struct {
	ATSURL string
}

func initConfig() {
	config = defaultConfig
	if url := os.Getenv("ATS_URL"); url != "" {
		if valid, _ := regexp.MatchString("https?://[a-zA-Z.0-9]+", strings.ToLower(url)); valid {
			config.ATSURL = url
		}

	}

}
