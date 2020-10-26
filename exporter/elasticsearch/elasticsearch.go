package elasticsearch

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/go-kit/kit/log"
	zapkit "github.com/go-kit/kit/log/zap"
	"github.com/justwatchcom/elasticsearch_exporter/collector"
	"github.com/justwatchcom/elasticsearch_exporter/pkg/clusterinfo"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
)

type Options struct {
	Timeout             time.Duration
	User                string
	Pass                string
	AllNodes            bool
	Indexes             bool
	Shards              bool
	Logger              *zap.SugaredLogger
	ClusterInfoInterval time.Duration
	Registry            *prometheus.Registry
}

var (
	logger log.Logger
)

func SetLogger(log *zap.SugaredLogger) {
	logger = zapkit.NewZapSugarLogger(log.Desugar(), zap.NewAtomicLevel().Level())
}

func DefaultOptions() Options {
	return Options{
		Timeout:             5 * time.Second,
		AllNodes:            true,
		Indexes:             false,
		Shards:              false,
		ClusterInfoInterval: 5 * time.Minute,
	}
}

func NewEsExporter(target string, opts Options) (err error) {
	httpClient := &http.Client{
		Timeout: opts.Timeout,
	}
	esURL, err := url.Parse(fmt.Sprintf("http://%s:%s@%s", opts.User, opts.Pass, target))
	if err != nil {
		return
	}
	clusterInfoRetriever := clusterinfo.New(logger, httpClient, esURL, 0)
	opts.Registry.MustRegister(collector.NewClusterHealth(logger, httpClient, esURL))
	opts.Registry.MustRegister(collector.NewNodes(logger, httpClient, esURL, opts.AllNodes, ""))
	if opts.Indexes || opts.Shards {
		iC := collector.NewIndices(logger, httpClient, esURL, opts.Shards)
		opts.Registry.MustRegister(iC)
		clusterInfoRetriever.RegisterConsumer(iC)
	}

	ctx, _ := context.WithCancel(context.Background())
	// start the cluster info retriever
	clusterInfoRetriever.Run(ctx)
	opts.Registry.MustRegister(clusterInfoRetriever)
	return
}
