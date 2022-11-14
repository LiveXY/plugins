package jsoner

import "github.com/livexy/plugins/plugin/streamer"

type JsonConfig struct {
	Path   string `yaml:"path"`   // 路径
	Driver string `yaml:"driver"` // 类型
}

// JSON接口
type Jsoner interface {
	streamer.Streamer
	MarshalString(val any) (string, error)
	UnmarshalString(buf string, val any) error
}
