package mysql

import (
	"context"
	"fmt"

	"github.com/go-kit/kit/log"
	zapkit "github.com/go-kit/kit/log/zap"
	"github.com/mxlxm/mysqld_exporter/collector"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
)

type Options struct {
	User     string
	Pass     string
	Scrapers map[collector.Scraper]bool
	Registry *prometheus.Registry
}

var scrapers = map[collector.Scraper]bool{
	collector.ScrapeGlobalStatus{}:                        true,
	collector.ScrapeGlobalVariables{}:                     true,
	collector.ScrapeSlaveStatus{}:                         true,
	collector.ScrapeProcesslist{}:                         false,
	collector.ScrapeUser{}:                                false,
	collector.ScrapeTableSchema{}:                         false,
	collector.ScrapeInfoSchemaInnodbTablespaces{}:         false,
	collector.ScrapeInnodbMetrics{}:                       false,
	collector.ScrapeAutoIncrementColumns{}:                false,
	collector.ScrapeBinlogSize{}:                          false,
	collector.ScrapePerfTableIOWaits{}:                    false,
	collector.ScrapePerfIndexIOWaits{}:                    false,
	collector.ScrapePerfTableLockWaits{}:                  false,
	collector.ScrapePerfEventsStatements{}:                false,
	collector.ScrapePerfEventsStatementsSum{}:             false,
	collector.ScrapePerfEventsWaits{}:                     false,
	collector.ScrapePerfFileEvents{}:                      false,
	collector.ScrapePerfFileInstances{}:                   false,
	collector.ScrapePerfReplicationGroupMemberStats{}:     false,
	collector.ScrapePerfReplicationApplierStatsByWorker{}: false,
	collector.ScrapeUserStat{}:                            false,
	collector.ScrapeClientStat{}:                          false,
	collector.ScrapeTableStat{}:                           false,
	collector.ScrapeSchemaStat{}:                          false,
	collector.ScrapeInnodbCmp{}:                           true,
	collector.ScrapeInnodbCmpMem{}:                        true,
	collector.ScrapeQueryResponseTime{}:                   true,
	collector.ScrapeEngineTokudbStatus{}:                  false,
	collector.ScrapeEngineInnodbStatus{}:                  false,
	collector.ScrapeHeartbeat{}:                           false,
	collector.ScrapeSlaveHosts{}:                          false,
}

func DefaultOptions() Options {
	return Options{
		Scrapers: scrapers,
	}
}

var (
	logger log.Logger
)

func SetLogger(l *zap.SugaredLogger) {
	logger = zapkit.NewZapSugarLogger(l.Desugar(), zap.NewAtomicLevel().Level())
}

func NewMysqlExporter(target string, opts Options) (err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/", opts.User, opts.Pass, target)
	ctx, _ := context.WithCancel(context.Background())
	enabledScrapers := []collector.Scraper{}
	for k, v := range opts.Scrapers {
		if v {
			enabledScrapers = append(enabledScrapers, k)
		}
	}
	opts.Registry.MustRegister(collector.New(ctx, dsn, collector.NewMetrics(), enabledScrapers, logger))
	return
}
