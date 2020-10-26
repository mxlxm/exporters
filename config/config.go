package config

import (
	"github.com/BurntSushi/toml"

	"github.com/mxlxm/common/utils"
)

type ExportersConfig struct {
	Listen     string
	ScrapePath string
	MetricPath string
	Exporters  map[string]Exporter
	Log        utils.LogConfig
}

type Auth struct {
	Instance string
	User     string
	Pass     string
}

type Exporter struct {
	DefaultUser string
	DefaultPass string
	Auths       []Auth
}

func LoadCFG(f string) (c ExportersConfig) {
	if _, err := toml.DecodeFile(f, &c); err != nil {
		panic(err)
	}
	return
}
