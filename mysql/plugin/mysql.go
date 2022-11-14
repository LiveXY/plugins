package plugin

import (
	"time"

	"github.com/livexy/plugins/plugin/dber"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"gorm.io/plugin/dbresolver"
)

type mysqlDb struct {
	db *gorm.DB
}

func New() dber.Dber {
	return &mysqlDb{}
}

func (p mysqlDb) Init(logname string, dbconf dber.DBConfig, logger logger.Interface) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(dbconf.Sources[0]), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
		Logger:                 logger,
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
		sources = append(sources, mysql.Open(v))
	}
	if len(sources) > 0 {
		conf.Sources = sources
	}
	replicas := []gorm.Dialector{}
	for _, v := range dbconf.Replicas {
		replicas = append(replicas, mysql.Open(v))
	}
	if len(replicas) > 0 {
		conf.Replicas = replicas
	}
	db.Use(dbresolver.Register(conf).
		SetMaxIdleConns(dbconf.MaxIdleConns).
		SetMaxOpenConns(dbconf.MaxOpenConns).
		SetConnMaxLifetime(time.Hour))
	return db, err
}

func (p mysqlDb) IfNull() string {
	return "ifnull"
}
func (p mysqlDb) If() string {
	return "if"
}
func (p mysqlDb) GroupConcat(field string) string {
	return "group_concat(" + field + ")"
}
func (p mysqlDb) GetSlots() int {
	return 65536
}
func (p mysqlDb) ClobScan(clob *dber.Clob, v any) error { return nil }

func (p mysqlDb) GetCreateID(value any, table, pk string) int64 {
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
