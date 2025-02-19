package opengauss

import (
	"database/sql"
	"fmt"
	"math"
	"strconv"
	"strings"

	_ "gitee.com/opengauss/openGauss-connector-go-pq"
	"gorm.io/gorm"
	"gorm.io/gorm/callbacks"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/migrator"
	"gorm.io/gorm/schema"
)

type Dialector struct {
	*Config
}

type Config struct {
	DriverName           string
	DSN                  string
	PreferSimpleProtocol bool
	WithoutReturning     bool
	Conn                 *sql.DB
}

func Open(dsn string) gorm.Dialector {
	return &Dialector{&Config{DSN: dsn}}
}

func New(config Config) gorm.Dialector {
	return &Dialector{Config: &config}
}

func (dialector Dialector) Name() string {
	return "opengauss"
}

func (dialector Dialector) Initialize(db *gorm.DB) (err error) {
	if dialector.DriverName == "" {
		dialector.DriverName = "opengauss"
	}
	if dialector.Conn != nil {
		db.ConnPool = dialector.Conn
	} else {
		db.ConnPool, err = sql.Open(dialector.DriverName, dialector.DSN)
		if err != nil {
			return err
		}
	}
	callbackConfig := &callbacks.Config{
		CreateClauses: []string{"INSERT", "VALUES", "ON CONFLICT"},
		QueryClauses:  []string{},
		UpdateClauses: []string{"UPDATE", "SET", "WHERE", "ORDER BY", "LIMIT"},
		DeleteClauses: []string{"DELETE", "FROM", "WHERE", "ORDER BY", "LIMIT"},
	}
	callbackConfig.LastInsertIDReversed = true
	callbackConfig.CreateClauses = append(callbackConfig.CreateClauses, "RETURNING")
	//callbackConfig.UpdateClauses = append(callbackConfig.UpdateClauses, "RETURNING")
	//callbackConfig.DeleteClauses = append(callbackConfig.DeleteClauses, "RETURNING")
	callbacks.RegisterDefaultCallbacks(db, callbackConfig)
	for k, v := range dialector.ClauseBuilders() {
		db.ClauseBuilders[k] = v
	}
	return
}

func (dialector Dialector) Migrator(db *gorm.DB) gorm.Migrator {
	return Migrator{migrator.Migrator{Config: migrator.Config{
		DB:                          db,
		Dialector:                   dialector,
		CreateIndexAfterCreateTable: true,
	}}}
}

func (dialector Dialector) DefaultValueOf(field *schema.Field) clause.Expression {
	return clause.Expr{SQL: "DEFAULT"}
}

func (dialector Dialector) BindVarTo(writer clause.Writer, stmt *gorm.Statement, v interface{}) {
	writer.WriteByte('$')
	writer.WriteString(strconv.Itoa(len(stmt.Vars)))
}

func (d Dialector) QuoteTo(writer clause.Writer, str string) {
	writer.WriteString(str)
}

func (dialector Dialector) Explain(sql string, vars ...interface{}) string {
	return logger.ExplainSQL(sql, nil, `'`, vars...)
}

func (dialector Dialector) DataTypeOf(field *schema.Field) string {
	switch field.DataType {
	case schema.Bool:
		return "boolean"
	case schema.Int, schema.Uint:
		return dialector.getSchemaIntAndUnitType(field)
	case schema.Float:
		return dialector.getSchemaFloatType(field)
	case schema.String:
		return dialector.getSchemaStringType(field)
	case schema.Time:
		return dialector.getSchemaTimeType(field)
	case schema.Bytes:
		return dialector.getSchemaBytesType(field)
	default:
		return dialector.getSchemaCustomType(field)
	}
}

func (dialector Dialector) getSchemaFloatType(field *schema.Field) string {
	if field.Precision > 0 {
		return fmt.Sprintf("decimal(%d, %d)", field.Precision, field.Scale)
	}
	if field.Size <= 32 {
		return "float"
	}
	return "double"
}

func (dialector Dialector) getSchemaStringType(field *schema.Field) string {
	size := field.Size
	if size == 0 {
		hasIndex := field.TagSettings["INDEX"] != "" || field.TagSettings["UNIQUE"] != ""
		if field.PrimaryKey || field.HasDefaultValue || hasIndex {
			size = 191 // utf8mb4
		}
	}
	if size >= 65536 && size <= int(math.Pow(2, 24)) {
		return "mediumtext"
	}
	if size > int(math.Pow(2, 24)) || size <= 0 {
		return "longtext"
	}
	return fmt.Sprintf("varchar(%d)", size)
}

func (dialector Dialector) getSchemaTimeType(field *schema.Field) string {
	var precision string
	if field.Precision > 0 {
		precision = fmt.Sprintf("(%d)", field.Precision)
	}

	if field.NotNull || field.PrimaryKey {
		return "datetime" + precision
	}
	return "datetime" + precision + " NULL"
}

func (dialector Dialector) getSchemaBytesType(field *schema.Field) string {
	if field.Size > 0 && field.Size < 65536 {
		return fmt.Sprintf("varbinary(%d)", field.Size)
	}

	if field.Size >= 65536 && field.Size <= int(math.Pow(2, 24)) {
		return "mediumblob"
	}

	return "longblob"
}

func (dialector Dialector) getSchemaIntAndUnitType(field *schema.Field) string {
	constraint := func(sqlType string) string {
		if field.DataType == schema.Uint {
			sqlType += " unsigned"
		}
		if field.AutoIncrement {
			sqlType += " AUTO_INCREMENT"
		}
		return sqlType
	}
	switch {
	case field.Size <= 8:
		return constraint("tinyint")
	case field.Size <= 16:
		return constraint("smallint")
	case field.Size <= 24:
		return constraint("mediumint")
	case field.Size <= 32:
		return constraint("int")
	default:
		return constraint("bigint")
	}
}

func (dialector Dialector) getSchemaCustomType(field *schema.Field) string {
	sqlType := string(field.DataType)
	if field.AutoIncrement && !strings.Contains(strings.ToLower(sqlType), " auto_increment") {
		sqlType += " AUTO_INCREMENT"
	}
	return sqlType
}

func (dialector Dialector) SavePoint(tx *gorm.DB, name string) error {
	return tx.Exec("SAVEPOINT " + name).Error
}

func (dialector Dialector) RollbackTo(tx *gorm.DB, name string) error {
	return tx.Exec("ROLLBACK TO SAVEPOINT " + name).Error
}

const (
	ClauseOnConflict = "ON CONFLICT"
	ClauseValues = "VALUES"
	ClauseFor = "FOR"
)

func (dialector Dialector) ClauseBuilders() map[string]clause.ClauseBuilder {
	clauseBuilders := map[string]clause.ClauseBuilder{
		ClauseOnConflict: func(c clause.Clause, builder clause.Builder) {
			onConflict, ok := c.Expression.(clause.OnConflict)
			if !ok {
				c.Build(builder)
				return
			}
			builder.WriteString("ON DUPLICATE KEY UPDATE ")
			if len(onConflict.DoUpdates) == 0 {
				if s := builder.(*gorm.Statement).Schema; s != nil {
					var column clause.Column
					onConflict.DoNothing = false
					if s.PrioritizedPrimaryField != nil {
						column = clause.Column{Name: s.PrioritizedPrimaryField.DBName}
					} else if len(s.DBNames) > 0 {
						column = clause.Column{Name: s.DBNames[0]}
					}
					if column.Name != "" {
						onConflict.DoUpdates = []clause.Assignment{{Column: column, Value: column}}
					}
					builder.(*gorm.Statement).AddClause(onConflict)
				}
			}
			for idx, assignment := range onConflict.DoUpdates {
				if idx > 0 {
					builder.WriteByte(',')
				}
				builder.WriteQuoted(assignment.Column)
				builder.WriteByte('=')
				if column, ok := assignment.Value.(clause.Column); ok && column.Table == "excluded" {
					column.Table = ""
					builder.WriteString("VALUES(")
					builder.WriteQuoted(column)
					builder.WriteByte(')')
				} else {
					builder.AddVar(builder, assignment.Value)
				}
			}
		},
		ClauseValues: func(c clause.Clause, builder clause.Builder) {
			if values, ok := c.Expression.(clause.Values); ok && len(values.Columns) == 0 {
				builder.WriteString("VALUES()")
				return
			}
			c.Build(builder)
		},
	}
	return clauseBuilders
}
