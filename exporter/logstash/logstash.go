package logstash

import (
	"fmt"

	"github.com/mxlxm/logstash_exporter/collector"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
)

var (
	log *zap.SugaredLogger
)

func SetLogger(l *zap.SugaredLogger) {
	log = l
}

type Options struct {
	Registry *prometheus.Registry
}

func DefaultOptions() Options {
	return Options{}
}

func NewLogstashExporter(target string, opts Options) (err error) {
	target = fmt.Sprintf("http://%s", target)
	stats, err := collector.NewNodeStatsCollector(target)
	if err != nil {
		log.Errorw("failed to get node status",
			"service", "logstash",
			"target", target,
			"error", err.Error(),
		)
		return
	}
	info, err := collector.NewNodeInfoCollector(target)
	if err != nil {
		log.Errorw("failed to get node info",
			"service", "logstash",
			"target", target,
			"error", err.Error(),
		)
		return
	}
	opts.Registry.MustRegister(stats)
	opts.Registry.MustRegister(info)
	return
}
