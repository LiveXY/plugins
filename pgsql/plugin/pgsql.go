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

func New() dber.Dber {
	return &pgsqlDb{}
}

func (p pgsqlDb) Init(logname string, dbconf dber.DBConfig, logger logger.Interface) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.New(postgres.Config{DSN: dbconf.Sources[0], PreferSimpleProtocol: true}), &gorm.Config{
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
	db.Use(dbresolver.Register(conf).
		SetMaxIdleConns(dbconf.MaxIdleConns).
		SetMaxOpenConns(dbconf.MaxOpenConns).
		SetConnMaxLifetime(time.Hour))
	return db, err
}

func (p pgsqlDb) IfNull() string {
	return "coalesce"
}
func (p pgsqlDb) If() string {
	return "iif"
}
func (p pgsqlDb) GroupConcat(field string) string {
	return "string_agg(" + field + ", ',')"
}
func (p pgsqlDb) GetSlots() int {
	return 65536
}
func (p pgsqlDb) ClobScan(clob *dber.Clob, v any) error { return nil }

func (p pgsqlDb) GetCreateID(value any, table, pk string) int64 {
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
