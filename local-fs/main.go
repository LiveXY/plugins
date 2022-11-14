package main

import (
	plug "github.com/livexy/plugins/local-fs/plugin"
	"github.com/livexy/plugins/plugin/osser"
)

var Plugin plugin
type plugin struct{}

func(p plugin) New(cfg osser.OSSConfig) (osser.OSSer, error) {
	return plug.NewLocalFS(cfg)
}
