package main

import (
	"github.com/livexy/plugin/exceler"
	plug "github.com/livexy/plugins/excel/plugin"
)

var Plugin plugin
type plugin struct{}

func(p plugin) New(cfg exceler.ExcelConfig) exceler.Exceler {
	return plug.NewQAXExcel(cfg)
}
