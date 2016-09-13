package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
)

var client = &http.Client{Timeout: 10 * time.Second}

func loadMetrics(config ATSExporterConfig) (*json.Decoder, error) {
	req, err := http.NewRequest("GET", config.ATSURL, nil)

	resp, err := client.Do(req)

	if err != nil || resp == nil || resp.StatusCode != 200 {
		status := 0
		if resp != nil {
			status = resp.StatusCode
		}
		log.WithFields(log.Fields{"error": err, "host": config.ATSURL, "statusCode": status}).Error("Error while retrieving data from ATSHost")
		return nil, errors.New("Error while retrieving data from ATSHost")
	}

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}
	return json.NewDecoder(bytes.NewBuffer(body)), nil
}

func getStatsInfo(config ATSExporterConfig) ([]StatsInfo, error) {
	var q []StatsInfo

	d, err := loadMetrics(config)
	if err != nil {
		return q, err
	}

	q = MakeStatsInfo(d)

	return q, nil
}

func getMetricMap(config ATSExporterConfig) (MetricMap, error) {
	var overview MetricMap

	d, err := loadMetrics(config)
	if err != nil {
		return overview, err
	}
	return MakeMap(d), nil
}
