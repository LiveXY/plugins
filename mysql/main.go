package main

import (
	plug "github.com/livexy/plugins/mysql/plugin"
	"github.com/livexy/plugins/plugin/dber"
)

var Plugin plugin

type plugin struct {}

func (p plugin) New() dber.Dber {
	return 	plug.New()
}
