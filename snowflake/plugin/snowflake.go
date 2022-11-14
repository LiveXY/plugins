package plugin

import (
	"sync"

	"github.com/livexy/plugin/ider"

	"github.com/bwmarrin/snowflake"
)
var (
	once sync.Once
	node *snowflake.Node
)

func New(cfg ider.IDConfig) string {
	once.Do(func() {
		idnode, err := snowflake.NewNode(int64(cfg.Node))
		if err != nil {
			panic("snowflake分布式ID生成算法配置错误：" + err.Error())
		}
		node = idnode
	})
	return node.Generate().String()
}
