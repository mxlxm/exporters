package exporter

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	failCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "exporters_scrape_failed",
			Help: "failed scrape counter ",
		}, []string{"service", "reason"},
	)
	invalidCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "exporters_invalid_scrape",
			Help: "failed scrape counter ",
		}, []string{"service", "target", "reason"},
	)
	scrapeDurationSummary = promauto.NewSummaryVec(
		prometheus.SummaryOpts{
			Name:       "exporters_scrape_duration_summary",
			Help:       "scrape duration summary",
			Objectives: map[float64]float64{0.5: 0.05, 0.75: 0.025, 0.95: 0.01, 0.99: 0.001},
		}, []string{"service", "target"},
	)
	MetricHandler = promhttp.Handler()
)
