package ider

type IDConfig struct {
	Path   string `yaml:"path"`   // 路径
	Driver string `yaml:"driver"` // ID类型
	Node   uint16 // 节点
}
