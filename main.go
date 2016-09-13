package main

import (
	"net/http"
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/prometheus/client_golang/prometheus"
)

const defaultLogLevel = log.InfoLevel

func initLogger() {
	log.SetLevel(getLogLevel())
	log.SetFormatter(&log.TextFormatter{})
}

func main() {
	initConfig()
	initLogger()
	exporter := newExporter()
	prometheus.MustRegister(exporter)

	http.Handle("/metrics", prometheus.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
             <head><title>ATS Exporter</title></head>
             <body>
             <h1>Apache Traffic Server Exporter</h1>
             <p><a href='/metrics'>Metrics</a></p>
             </body>
             </html>`))
	})

	log.WithFields(log.Fields{
		"ATS_URL":    config.ATSURL,
		"VERSION":    Version,
		"REVISION":   Revision,
		"BRANCH":     Branch,
		"BUILD_DATE": BuildDate,
	}).Info("Starting ATS exporter")

	log.Fatal(http.ListenAndServe(":9090", nil))
}

func getLogLevel() log.Level {
	lvl := strings.ToLower(os.Getenv("LOG_LEVEL"))
	level, err := log.ParseLevel(lvl)
	if err != nil {
		level = defaultLogLevel
	}
	return level
}
