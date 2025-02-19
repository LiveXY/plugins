package plugin

import (
	"time"

	"github.com/livexy/plugin/dber"
	"github.com/livexy/plugins/opengaussb/openguass"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"gorm.io/plugin/dbresolver"
)

type gaussDb struct {
	db *gorm.DB
}

func New() dber.Dber {
	return &gaussDb{}
}

func (p gaussDb) Init(logname string, dbconf dber.DBConfig, logger logger.Interface) (*gorm.DB, error) {
	db, err := gorm.Open(openguass.New(openguass.Config{DSN: dbconf.Sources[0]}), &gorm.Config{
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
		sources = append(sources, openguass.New(openguass.Config{DSN: v}))
	}
	if len(sources) > 0 {
		conf.Sources = sources
	}
	replicas := []gorm.Dialector{}
	for _, v := range dbconf.Replicas {
		replicas = append(replicas, openguass.New(openguass.Config{DSN: v}))
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
func (p gaussDb) ExAdd(field string, val any) clause.Expr {
	return gorm.Expr(field + `+?`, val)
}
func (p gaussDb) IfNull() string {
	return "ifnull"
}
func (p gaussDb) If() string {
	return "if"
}
func (p gaussDb) GroupConcat(field string) string {
	return "group_concat(" + field + ")"
}
func (p gaussDb) GetSlots() int {
	return 65536
}
func (p gaussDb) ClobScan(clob *dber.Clob, v any) error { return nil }

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
