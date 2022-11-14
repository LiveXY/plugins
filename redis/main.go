package main

import (
	"github.com/livexy/plugins/plugin/cacher"
	plug "github.com/livexy/plugins/redis/plugin"

	"go.uber.org/zap"
)

var Plugin plugin
type plugin struct{}

func(p plugin) New(cfg cacher.CacheConfig, logger *zap.Logger) (cacher.Cacher, error) {
	return plug.NewRedisCache(cfg, logger)
}
