package exporter

import (
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/mxlxm/exporters/exporter/elasticsearch"
	"github.com/mxlxm/exporters/exporter/kafka"
	"github.com/mxlxm/exporters/exporter/logstash"
	"github.com/mxlxm/exporters/exporter/mysql"
	"github.com/mxlxm/exporters/exporter/redis"
	"github.com/mxlxm/exporters/exporter/rocketmq"
)

func ScrapeHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	service := strings.ToLower(r.URL.Query().Get("service"))
	if service == "" {
		http.Error(w, "'service' parameter must be specified", 400)
		invalidCounter.WithLabelValues("", "", "no service").Inc()
		return
	}
	target := r.URL.Query().Get("target")
	if target == "" {
		http.Error(w, "'target' parameter must be specified", 400)
		invalidCounter.WithLabelValues(service, "", "no target").Inc()
		return
	}
	registry, new := GetOrNewRegistry(target)
	if new {
		var err error
		user, pass := getAuth(service, target)
		switch service {
		case "redis":
			opts := redis.DefaultOptions()
			opts.Registry = registry
			opts.Password = pass
			redis.SetLogger(log)
			err = redis.NewRedisExporter(target, opts)
		case "elasticsearch":
			opts := elasticsearch.DefaultOptions()
			opts.Registry = registry
			opts.User = user
			opts.Pass = pass
			elasticsearch.SetLogger(log)
			err = elasticsearch.NewEsExporter(target, opts)
		case "mysql", "polardb":
			opts := mysql.DefaultOptions()
			opts.Registry = registry
			opts.User = user
			opts.Pass = pass
			mysql.SetLogger(log)
			err = mysql.NewMysqlExporter(target, opts)
		case "kafka":
			// get cluster name
			cluster := strings.ToLower(r.URL.Query().Get("cluster"))
			if cluster == "" {
				http.Error(w, "'cluster' parameter must be specified for kafka", 400)
				invalidCounter.WithLabelValues(service, target, "no cluster").Inc()
				return
			}
			// cache cluster nodes
			cnodes := kafka.AddClusterNode(cluster, target)
			// return if there's no enough nodes
			if len(cnodes) < 3 {
				TargetRegistry.Remove(target)
				return
			}
			nodes := make([]string, 0, len(cnodes))
			for k := range cnodes {
				nodes = append(nodes, k)
			}
			sort.Slice(nodes, func(i, j int) bool { return nodes[i] < nodes[j] })
			// only response for first node, cluster level metrics are the same.
			if target != nodes[0] {
				TargetRegistry.Remove(target)
				return
			}
			opts := kafka.DefaultOptions()
			opts.Registry = registry
			kafka.SetLogger(log)
			err = kafka.NewKafkaExporter(nodes, opts)
		case "rocketmq":
			opts := rocketmq.DefaultOptions()
			opts.Registry = registry
			rocketmq.SetLogger(log)
			err = rocketmq.NewRocketMQExporter(target, opts)
		case "logstash":
			opts := logstash.DefaultOptions()
			opts.Registry = registry
			logstash.SetLogger(log)
			err = logstash.NewLogstashExporter(target, opts)
		default:
			http.Error(w, fmt.Sprintf("unsupported service:%s", service), 503)
			invalidCounter.WithLabelValues(service, target, "unsupported").Inc()
			return
		}
		if err != nil {
			http.Error(w, "failed to new exporter", 500)
			TargetRegistry.Remove(target)
			failCounter.WithLabelValues(service, target, err.Error()).Inc()
			return
		}
	}
	promhttp.HandlerFor(
		registry, promhttp.HandlerOpts{ErrorHandling: promhttp.ContinueOnError},
	).ServeHTTP(w, r)
	cost := time.Now().Sub(start)
	scrapeDurationSummary.WithLabelValues(service, target).Observe(float64(cost.Milliseconds()))
}
