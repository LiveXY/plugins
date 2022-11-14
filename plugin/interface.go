package plugin

import (
	"github.com/livexy/plugins/plugin/cacher"
	"github.com/livexy/plugins/plugin/dber"
	"github.com/livexy/plugins/plugin/exceler"
	"github.com/livexy/plugins/plugin/ider"
	"github.com/livexy/plugins/plugin/jsoner"
	"github.com/livexy/plugins/plugin/osser"
	"github.com/livexy/plugins/plugin/smser"
	"github.com/livexy/plugins/plugin/streamer"
	"github.com/livexy/plugins/plugin/worder"

	"go.uber.org/zap"
)

type DBer interface {
	New() dber.Dber
}

type IDer interface {
	New(cfg ider.IDConfig) string
}

type Cacher interface {
	New(cfg cacher.CacheConfig, logger *zap.Logger) (cacher.Cacher, error)
}

type OSSer interface {
	New(cfg osser.OSSConfig) (osser.OSSer, error)
}

type Exceler interface {
	New(cfg exceler.ExcelConfig) exceler.Exceler
}

type Worder interface {
	New(cfg worder.WordConfig) worder.Worder
}

type Jsoner interface {
	New(cfg jsoner.JsonConfig) jsoner.Jsoner
}

type Streamer interface {
	New(cfg streamer.StreamConfig) streamer.Streamer
}

type SMSer interface {
	New(cfg smser.SMSConfig) smser.SMSer
}
