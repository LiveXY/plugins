package main

import (
	"github.com/livexy/plugin/dber"
	plug "github.com/livexy/plugins/opengaussb/plugin"
)

var Plugin plugin

type plugin struct{}

func (p plugin) New() dber.Dber {
	return plug.New()
}
