package streamer

type StreamConfig struct {
	Path   string `yaml:"path"`   // 路径
	Driver string `yaml:"driver"` // 类型
}

// GOB接口
type Streamer interface {
	Marshal(val any) ([]byte, error)
	Unmarshal(buf []byte, val any) error
}
