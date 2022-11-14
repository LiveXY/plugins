package main

import (
	"github.com/livexy/plugin/osser"
	plug "github.com/livexy/plugins/local-fs/plugin"
)

var Plugin plugin
type plugin struct{}

func(p plugin) New(cfg osser.OSSConfig) (osser.OSSer, error) {
	return plug.NewLocalFS(cfg)
}
