package main

import (
	"github.com/livexy/plugin/smser"
	plug "github.com/livexy/plugins/jianzhou-sms/plugin"
)

var Plugin plugin
type plugin struct{}

func(p plugin) New(cfg smser.SMSConfig) smser.SMSer {
	return plug.NewSMS(cfg)
}
