package main

import (
	"github.com/livexy/plugin/jsoner"
	plug "github.com/livexy/plugins/go-json/plugin"
)

var Plugin plugin
type plugin struct{}

func(p plugin) New(cfg jsoner.JsonConfig) jsoner.Jsoner {
	return plug.NewGoJson(cfg)
}
