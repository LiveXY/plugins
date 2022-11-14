package main

import (
	"github.com/livexy/plugins/plugin/ider"

	"github.com/rs/xid"
)

var Plugin plugin
type plugin struct{}

func(p plugin) New(cfg ider.IDConfig) string {
	return xid.New().String()
}
