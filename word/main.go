package main

import (
	"github.com/livexy/plugin/worder"
	plug "github.com/livexy/plugins/word/plugin"
)

var Plugin plugin
type plugin struct{}

func(p plugin) New(cfg worder.WordConfig) worder.Worder {
	return plug.NewDocxWord(cfg)
}
