package main

import (
	plug "github.com/livexy/plugins/jianzhou-sms/plugin"
	"github.com/livexy/plugins/plugin/smser"
)

var Plugin plugin
type plugin struct{}

func(p plugin) New(cfg smser.SMSConfig) smser.SMSer {
	return plug.NewSMS(cfg)
}
