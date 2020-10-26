package exporter

import (
	"github.com/orcaman/concurrent-map"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"

	"github.com/mxlxm/exporters/config"
)

var (
	log            *zap.SugaredLogger
	InstanceAuths  = cmap.New()
	ServiceAuths   = cmap.New()
	TargetRegistry = cmap.New()
)

func SetLogger(l *zap.SugaredLogger) {
	log = l
}

func BuildAuth(cfg config.ExportersConfig) {
	for s, e := range cfg.Exporters {
		for _, a := range e.Auths {
			InstanceAuths.Set(a.Instance, a)
		}
		ServiceAuths.Set(s, e)
	}
}

func getAuth(service, target string) (user, pass string) {
	if InstanceAuths.Has(target) {
		value, _ := InstanceAuths.Get(target)
		auth := value.(config.Auth)
		return auth.User, auth.Pass
	}
	if ServiceAuths.Has(service) {
		value, _ := ServiceAuths.Get(service)
		auth := value.(config.Exporter)
		return auth.DefaultUser, auth.DefaultPass
	}
	return
}

func GetOrNewRegistry(target string) (registry *prometheus.Registry, new bool) {
	if v, exist := TargetRegistry.Get(target); exist {
		registry = v.(*prometheus.Registry)
		new = false
	} else {
		registry = prometheus.NewRegistry()
		TargetRegistry.Set(target, registry)
		new = true
	}
	return
}
