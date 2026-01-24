package plugin

import (
	"time"

	"github.com/livexy/plugin/dber"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"gorm.io/plugin/dbresolver"
)

type pgsqlDb struct {
	db *gorm.DB
}

// New 创建一个新的 PostgreSQL 数据库适配实例
func New() dber.Dber {
	return &pgsqlDb{}
}

// Init 初始化数据库连接
func (p pgsqlDb) Init(logname string, dbconf dber.DBConfig, val any) (any, error) {
	var l logger.Interface
	if v, ok := val.(logger.Interface); ok {
		l = v
	}
	db, err := gorm.Open(postgres.New(postgres.Config{DSN: dbconf.Sources[0], PreferSimpleProtocol: true}), &gorm.Config{
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
		sources = append(sources, postgres.New(postgres.Config{DSN: v, PreferSimpleProtocol: true}))
	}
	if len(sources) > 0 {
		conf.Sources = sources
	}
	replicas := []gorm.Dialector{}
	for _, v := range dbconf.Replicas {
		replicas = append(replicas, postgres.New(postgres.Config{DSN: v, PreferSimpleProtocol: true}))
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

// ExAdd 返回用于 GORM 的字段自增表达式（针对 PostgreSQL 的 excluded 语法）
func (p pgsqlDb) ExAdd(field string, val any) any {
	return gorm.Expr(`"excluded"."`+field+`"+?`, val)
}

// IfNull 返回 PostgreSQL 的空值判断函数名 (coalesce)
func (p pgsqlDb) IfNull() string {
	return "coalesce"
}

// If 返回 PostgreSQL 的条件判断函数名 (iif)
func (p pgsqlDb) If() string {
	return "iif"
}

// GroupConcat 返回 PostgreSQL 的字符串聚合函数表达式
func (p pgsqlDb) GroupConcat(field string) string {
	return "string_agg(" + field + ", ',')"
}

// GetSlots 获取数据库支持的最大插槽数
func (p pgsqlDb) GetSlots() int {
	return 65536
}

// ClobScan 处理 CLOB 类型的扫描（PostgreSQL 默认返回 nil）
func (p pgsqlDb) ClobScan(clob *dber.Clob, v any) error { return nil }

// GetCreateID 插入数据并获取自增 ID
func (p pgsqlDb) GetCreateID(value any, table, pk string) int64 {
	tx := p.db.Begin()
	id := getInsertID(tx.Create(value), table, pk)
	tx.Commit()
	return id
}
func getInsertID(db *gorm.DB, table, pk string) int64 {
	var id int64
	db.Raw("select lastval() as id").Scan(&id)
	return id
}
