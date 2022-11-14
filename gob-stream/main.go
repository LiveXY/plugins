package main

import (
	"github.com/livexy/plugin/streamer"
	plug "github.com/livexy/plugins/gob-stream/plugin"
)

var Plugin plugin
type plugin struct{}

func(p plugin) New(cfg streamer.StreamConfig) streamer.Streamer {
	return plug.NewStream(cfg)
}
