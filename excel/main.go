package main

import (
	plug "github.com/livexy/plugins/excel/plugin"

	"github.com/livexy/plugin/exceler"
)

var Plugin plugin
type plugin struct{}

func(p plugin) New(cfg exceler.ExcelConfig) exceler.Exceler {
	return plug.NewQAXExcel(cfg)
}
