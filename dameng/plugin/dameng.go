package plugin

import (
	"time"

	"github.com/livexy/plugin/dber"
	"github.com/livexy/plugins/dameng/dameng"

	"gitee.com/chunanyong/dm"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"gorm.io/plugin/dbresolver"
)

type damengDb struct {
	db *gorm.DB
}

// New 创建一个新的 Dameng 数据库适配实例
func New() dber.Dber {
	return &damengDb{}
}

// Init 初始化数据库连接
// 包含连接池配置、读写分离配置及驱动特定的优化选项
func (p damengDb) Init(logname string, dbconf dber.DBConfig, val any) (any, error) {
	var l logger.Interface
	if v, ok := val.(logger.Interface); ok {
		l = v
	}
	db, err := gorm.Open(dameng.Open(dbconf.Sources[0]), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
		Logger:                 l,
	})
	if err != nil {
		return nil, err
	}
	conf := dbresolver.Config{Policy: dbresolver.RandomPolicy{}}
	sources := []gorm.Dialector{}
	for k, v := range dbconf.Sources {
		if k == 0 {
			continue
		}
		sources = append(sources, dameng.Open(v))
	}
	if len(sources) > 0 {
		conf.Sources = sources
	}
	replicas := []gorm.Dialector{}
	for _, v := range dbconf.Replicas {
		replicas = append(replicas, dameng.Open(v))
	}
	if len(replicas) > 0 {
		conf.Replicas = replicas
	}
	err = db.Use(dbresolver.Register(conf).
		SetMaxIdleConns(dbconf.MaxIdleConns).
		SetMaxOpenConns(dbconf.MaxOpenConns).
		SetConnMaxLifetime(time.Hour))
	return db, err
}

// IfNull 返回 Dameng 的空值判读函数名 (nvl)
func (p damengDb) IfNull() string {
	return "nvl"
}

// If 返回 Dameng 的条件判断函数名 (if)
func (p damengDb) If() string {
	return "if"
}

// GroupConcat 返回 Dameng 的字符串聚合函数表达式
func (p damengDb) GroupConcat(field string) string {
	return "wm_concat(" + field + ")"
}

// GetSlots 获取数据库支持的最大插槽数
func (p damengDb) GetSlots() int {
	return 32768
}

// ExAdd 返回用于 GORM 的字段自增表达式
func (p damengDb) ExAdd(field string, val any) any {
	return gorm.Expr(field+`+?`, val)
}

// ClobScan 处理 Dameng 数据库的 CLOB 类型扫描与转换
func (p damengDb) ClobScan(clob *dber.Clob, v any) error {
	switch val := v.(type) {
	case *dm.DmClob:
		le, err := val.GetLength()
		if err != nil {
			return err
		}
		if le == 0 {
			*clob = ""
		} else {
			str, err := val.ReadString(1, int(le))
			if err != nil {
				return nil
			}
			*clob = dber.Clob(str)
		}
	case []uint8:
		bs := []byte{}
		for _, b := range val {
			bs = append(bs, byte(b))
		}
		*clob = dber.Clob(string(bs))
	case string:
		*clob = dber.Clob(val)
	}
	return nil
}

// GetCreateID 插入数据并获取自增 ID
func (p damengDb) GetCreateID(value any, table, pk string) int64 {
	tx := p.db.Begin()
	id := getInsertID(tx.Create(value), table, pk)
	tx.Commit()
	return id
}
func getInsertID(db *gorm.DB, table, pk string) int64 {
	var id int64
	db.Raw("SELECT @@IDENTITY as id").Scan(&id)
	return id
}
