package main

import (
	plug "github.com/livexy/plugins/gob-stream/plugin"
	"github.com/livexy/plugins/plugin/streamer"
)

var Plugin plugin
type plugin struct{}

func(p plugin) New(cfg streamer.StreamConfig) streamer.Streamer {
	return plug.NewStream(cfg)
}
