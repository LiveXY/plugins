package plugin

import (
	"time"

	"github.com/livexy/plugin/dber"
	"github.com/livexy/plugins/dameng/dameng"

	"gitee.com/chunanyong/dm"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"gorm.io/plugin/dbresolver"
)

type damengDb struct {
	db *gorm.DB
}

func New() dber.Dber {
	return &damengDb{}
}

func (p damengDb) Init(logname string, dbconf dber.DBConfig, logger logger.Interface) (*gorm.DB, error) {
	db, err := gorm.Open(dameng.Open(dbconf.Sources[0]), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		QueryFields:                              true,
		DisableForeignKeyConstraintWhenMigrating: true,
		SkipDefaultTransaction:                   true,
		PrepareStmt:                              true,
		Logger:                                   logger,
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
	db.Use(dbresolver.Register(conf).
		SetMaxIdleConns(dbconf.MaxIdleConns).
		SetMaxOpenConns(dbconf.MaxOpenConns).
		SetConnMaxLifetime(time.Hour))
	return db, err
}

func (p damengDb) IfNull() string {
	return "nvl"
}
func (p damengDb) If() string {
	return "if"
}
func (p damengDb) GroupConcat(field string) string {
	return "wm_concat(" + field + ")"
}
func (p damengDb) GetSlots() int {
	return 32768
}
func (p damengDb) ExAdd(field string, val any) clause.Expr {
	return gorm.Expr(field + `+?`, val)
}
func (p damengDb) ClobScan(clob *dber.Clob, v any) error {
	switch v.(type) {
	case *dm.DmClob:
		tmp := v.(*dm.DmClob)
		le, err := tmp.GetLength()
		if err != nil {
			return err
		}
		if le == 0 {
			*clob = ""
		} else {
			str, err := tmp.ReadString(1, int(le))
			if err != nil {
				return nil
			}
			*clob = dber.Clob(str)
		}
	case []uint8:
		bs := []byte{}
		for _, b := range v.([]uint8) {
			bs = append(bs, byte(b))
		}
		*clob = dber.Clob(string(bs))
	case string:
		*clob = dber.Clob(v.(string))
	}
	return nil
}
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
