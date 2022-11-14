package dber

import (
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DBConfig struct {
	Path         string   `yaml:"path"` // 路径
	Driver       string   `yaml:"driver"`
	Sources      []string `yaml:"sources"`
	Replicas     []string `yaml:"replicas"`
	MaxIdleConns int      `yaml:"maxIdleConns"`
	MaxOpenConns int      `yaml:"maxOpenConns"`
}

type Dber interface {
	Init(logname string, cfg DBConfig, logger logger.Interface) (*gorm.DB, error)
	IfNull() string
	If() string
	GroupConcat(field string) string
	GetSlots() int
	GetCreateID(value any, table, pk string) int64
	ClobScan(clob *Clob, v any) error
}
