package main

import (
	"github.com/livexy/plugins/plugin/jsoner"
	plug "github.com/livexy/plugins/sonic-json/plugin"
)

var Plugin plugin
type plugin struct{}

func(p plugin) New(cfg jsoner.JsonConfig) jsoner.Jsoner {
	return plug.NewSonicJson(cfg)
}
