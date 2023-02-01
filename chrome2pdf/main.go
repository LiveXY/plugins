package main

import (
	"github.com/livexy/plugin/pdfer"
	plug "github.com/livexy/plugins/chrome2pdf/plugin"
)

var Plugin plugin
type plugin struct{}

func(p plugin) New() pdfer.PDFer {
	return plug.NewChromePdf()
}
