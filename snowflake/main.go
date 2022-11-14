package main

import (
	"github.com/livexy/plugins/plugin/ider"
	plug "github.com/livexy/plugins/snowflake/plugin"
)

var Plugin plugin
type plugin struct{}

func(p plugin) New(cfg ider.IDConfig) string {
	return plug.New(cfg)
}
