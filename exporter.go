package main

import (
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/prometheus/client_golang/prometheus"
)

type exporter struct {
	mutex    sync.RWMutex
	metrics  map[string]prometheus.Gauge
	upMetric prometheus.Gauge
}

func newExporter() *exporter {
	return &exporter{
		metrics:  metricDescription,
		upMetric: upMetricDescription,
	}
}

func (e *exporter) fetchATS() {
	ATSData, err := getMetricMap(config)

	if err != nil {
		e.upMetric.Set(0)
	} else {
		e.upMetric.Set(1)
	}

	for key, gauge := range e.metrics {
		if value, ok := ATSData[key]; ok {
			log.WithFields(log.Fields{"key": key, "value": value}).Debug("Set metric for key")
			gauge.Set(value)
		} else {
			log.WithFields(log.Fields{"key": key}).Warn("data not found")
		}
	}

	log.Info("Metrics updated successfully.")
}

func (e *exporter) Describe(ch chan<- *prometheus.Desc) {
	for _, gauge := range e.metrics {
		gauge.Describe(ch)
	}
}

func (e *exporter) Collect(ch chan<- prometheus.Metric) {
	e.mutex.Lock() // To protect metrics from concurrent collects.
	defer e.mutex.Unlock()

	e.fetchATS()

	e.upMetric.Collect(ch)

	for _, gauge := range e.metrics {
		gauge.Collect(ch)
	}

	BuildInfo.Collect(ch)
}
