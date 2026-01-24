package plugin

import (
	"time"

	"github.com/livexy/plugin/dber"
	"github.com/livexy/plugins/opengaussb/opengauss"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"gorm.io/plugin/dbresolver"
)

type gaussDb struct {
	db *gorm.DB
}

// New 创建一个新的 OpenGauss 数据库适配实例
func New() dber.Dber {
	return &gaussDb{}
}

// Init 初始化数据库连接
func (p gaussDb) Init(logname string, dbconf dber.DBConfig, val any) (any, error) {
	var l logger.Interface
	if v, ok := val.(logger.Interface); ok {
		l = v
	}
	db, err := gorm.Open(opengauss.New(opengauss.Config{DSN: dbconf.Sources[0]}), &gorm.Config{
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
		sources = append(sources, opengauss.New(opengauss.Config{DSN: v}))
	}
	if len(sources) > 0 {
		conf.Sources = sources
	}
	replicas := []gorm.Dialector{}
	for _, v := range dbconf.Replicas {
		replicas = append(replicas, opengauss.New(opengauss.Config{DSN: v}))
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

// ExAdd 返回用于 GORM 的字段自增表达式
func (p gaussDb) ExAdd(field string, val any) any {
	return gorm.Expr(field+`+?`, val)
}

// IfNull 返回 OpenGauss 的 IFNULL 函数名 (ifnull)
func (p gaussDb) IfNull() string {
	return "ifnull"
}

// If 返回 OpenGauss 的 IF 函数名 (if)
func (p gaussDb) If() string {
	return "if"
}

// GroupConcat 返回 OpenGauss 的 GroupConcat 表达式
func (p gaussDb) GroupConcat(field string) string {
	return "group_concat(" + field + ")"
}

// GetSlots 获取数据库支持的最大插槽数
func (p gaussDb) GetSlots() int {
	return 65536
}

// ClobScan 处理 CLOB 类型的扫描（OpenGauss 默认返回 nil）
func (p gaussDb) ClobScan(clob *dber.Clob, v any) error { return nil }

// GetCreateID 插入数据并获取自增 ID
func (p gaussDb) GetCreateID(value any, table, pk string) int64 {
	tx := p.db.Begin()
	id := getInsertID(tx.Create(value), table, pk)
	tx.Commit()
	return id
}
func getInsertID(db *gorm.DB, table, pk string) int64 {
	var id int64
	db.Raw("select LAST_INSERT_ID() as id").Scan(&id)
	return id
}
