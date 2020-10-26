package main

import (
	"flag"
	"net/http"
	"net/http/pprof"
	"path"

	"go.uber.org/zap"

	"github.com/mxlxm/common/utils"
	"github.com/mxlxm/exporters/config"
	"github.com/mxlxm/exporters/exporter"
)

var (
	err  error
	cfg  config.ExportersConfig
	log  *zap.SugaredLogger
	mux  *http.ServeMux
	conf = flag.String("conf", path.Clean(utils.ExecDir()+"/../conf/exporters.toml"), "exporters config file")
)

func init() {
	flag.Parse()
	cfg = config.LoadCFG(*conf)
	if log, err = utils.SugarInit(utils.InitLogConfig(cfg.Log)); err != nil {
		panic(err)
	}
	exporter.SetLogger(log)
	exporter.BuildAuth(cfg)
}

func main() {
	defer func() {
		log.Sync()
	}()

	mux = http.NewServeMux()

	mux.HandleFunc(cfg.ScrapePath, exporter.ScrapeHandler)
	mux.Handle(cfg.MetricPath, exporter.MetricHandler)

	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	http.ListenAndServe(cfg.Listen, mux)
}
