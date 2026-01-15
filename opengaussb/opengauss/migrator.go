package opengauss

import (
	"database/sql"
	"fmt"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/migrator"
	"gorm.io/gorm/schema"
)
const indexSql = `
SELECT
	TABLE_NAME,
	COLUMN_NAME,
	INDEX_NAME,
	NON_UNIQUE 
FROM
	information_schema.STATISTICS 
WHERE
	TABLE_SCHEMA = ? 
	AND TABLE_NAME = ? 
ORDER BY
	INDEX_NAME,
	SEQ_IN_INDEX`

var typeAliasMap = map[string][]string{
	"bool":    {"tinyint"},
	"tinyint": {"bool"},
}

type Migrator struct {
	migrator.Migrator
}

func (m Migrator) FullDataTypeOf(field *schema.Field) clause.Expr {
	expr := m.Migrator.FullDataTypeOf(field)
	if value, ok := field.TagSettings["COMMENT"]; ok {
		expr.SQL += " COMMENT " + m.Dialector.Explain("?", value)
	}
	return expr
}

func (m Migrator) AlterColumn(value interface{}, field string) error {
	return m.RunWithValue(value, func(stmt *gorm.Statement) error {
		if stmt.Schema != nil {
			if field := stmt.Schema.LookUpField(field); field != nil {
				fullDataType := m.FullDataTypeOf(field)
				return m.DB.Exec(
					"ALTER TABLE ? MODIFY COLUMN ? ?",
					clause.Table{Name: stmt.Table}, clause.Column{Name: field.DBName}, fullDataType,
				).Error
			}
		}
		return fmt.Errorf("failed to look up field with name: %s", field)
	})
}

func (m Migrator) CurrentDatabase() (name string) {
	baseName := m.Migrator.CurrentDatabase()
	m.DB.Raw(
		"SELECT SCHEMA_NAME from Information_schema.SCHEMATA where SCHEMA_NAME LIKE ? ORDER BY SCHEMA_NAME=? DESC,SCHEMA_NAME limit 1",
		baseName+"%", baseName).Scan(&name)
	return
}

func (m Migrator) GetTables() (tableList []string, err error) {
	err = m.DB.Raw("SELECT TABLE_NAME FROM information_schema.tables where TABLE_SCHEMA=?", m.CurrentDatabase()).
		Scan(&tableList).Error
	return
}

func (m Migrator) GetIndexes(value interface{}) ([]gorm.Index, error) {
	indexes := make([]gorm.Index, 0)
	err := m.RunWithValue(value, func(stmt *gorm.Statement) error {
		result := make([]*Index, 0)
		schema, table := m.CurrentSchema(stmt, stmt.Table)
		scanErr := m.DB.Table(table).Raw(indexSql, schema, table).Scan(&result).Error
		if scanErr != nil {
			return scanErr
		}
		indexMap, indexNames := groupByIndexName(result)
		for _, name := range indexNames {
			idx := indexMap[name]
			if len(idx) == 0 {
				continue
			}
			tempIdx := &migrator.Index{
				TableName: idx[0].TableName,
				NameValue: idx[0].IndexName,
				PrimaryKeyValue: sql.NullBool{
					Bool:  idx[0].IndexName == "PRIMARY",
					Valid: true,
				},
				UniqueValue: sql.NullBool{
					Bool:  idx[0].NonUnique == 0,
					Valid: true,
				},
			}
			for _, x := range idx {
				tempIdx.ColumnList = append(tempIdx.ColumnList, x.ColumnName)
			}
			indexes = append(indexes, tempIdx)
		}
		return nil
	})
	return indexes, err
}

func (m Migrator) BuildIndexOptions(opts []schema.IndexOption, stmt *gorm.Statement) (results []interface{}) {
	for _, opt := range opts {
		str := stmt.Quote(opt.DBName)
		if opt.Expression != "" {
			str = opt.Expression
		}

		if opt.Collate != "" {
			str += " COLLATE " + opt.Collate
		}

		if opt.Sort != "" {
			str += " " + opt.Sort
		}
		results = append(results, clause.Expr{SQL: str})
	}
	return
}

func (m Migrator) HasIndex(value interface{}, name string) bool {
	var count int64
	_ = m.RunWithValue(value, func(stmt *gorm.Statement) error {
		if stmt.Schema != nil {
			if idx := stmt.Schema.LookIndex(name); idx != nil {
				name = idx.Name
			}
		}
		currentSchema, curTable := m.CurrentSchema(stmt, stmt.Table)
		return m.DB.Raw(
			"SELECT count(*) FROM pg_indexes WHERE tablename = ? AND indexname = ? AND schemaname = ?", curTable, name, currentSchema,
		).Scan(&count).Error
	})

	return count > 0
}

func (m Migrator) CreateIndex(value interface{}, name string) error {
	return m.RunWithValue(value, func(stmt *gorm.Statement) error {
		if stmt.Schema != nil {
			if idx := stmt.Schema.LookIndex(name); idx != nil {
				opts := m.BuildIndexOptions(idx.Fields, stmt)
				values := []interface{}{clause.Column{Name: idx.Name}, m.CurrentTable(stmt), opts}

				createIndexSQL := "CREATE "
				if idx.Class != "" {
					createIndexSQL += idx.Class + " "
				}
				createIndexSQL += "INDEX "

				if strings.TrimSpace(strings.ToUpper(idx.Option)) == "CONCURRENTLY" {
					createIndexSQL += "CONCURRENTLY "
				}

				createIndexSQL += "IF NOT EXISTS ? ON ?"

				if idx.Type != "" {
					createIndexSQL += " USING " + idx.Type + "(?)"
				} else {
					createIndexSQL += " ?"
				}

				if idx.Where != "" {
					createIndexSQL += " WHERE " + idx.Where
				}

				return m.DB.Exec(createIndexSQL, values...).Error
			}
		}

		return fmt.Errorf("failed to create index with name %v", name)
	})
}

func (m Migrator) RenameIndex(value interface{}, oldName, newName string) error {
	
	return m.RunWithValue(value, func(stmt *gorm.Statement) error {
		err := m.DropIndex(value, oldName)
		if err != nil {
			return err
		}

		if stmt.Schema != nil {
			if idx := stmt.Schema.LookIndex(newName); idx == nil {
				if idx = stmt.Schema.LookIndex(oldName); idx != nil {
					opts := m.BuildIndexOptions(idx.Fields, stmt)
					values := []interface{}{clause.Column{Name: newName}, clause.Table{Name: stmt.Table}, opts}

					createIndexSQL := "CREATE "
					if idx.Class != "" {
						createIndexSQL += idx.Class + " "
					}
					createIndexSQL += "INDEX ? ON ??"

					if idx.Type != "" {
						createIndexSQL += " USING " + idx.Type
					}

					return m.DB.Exec(createIndexSQL, values...).Error
				}
			}
		}

		return m.CreateIndex(value, newName)
	})

}
func (m Migrator) DropIndex(value interface{}, name string) error {
	return m.RunWithValue(value, func(stmt *gorm.Statement) error {
		if stmt.Schema != nil {
			if idx := stmt.Schema.LookIndex(name); idx != nil {
				name = idx.Name
			}
		}

		return m.DB.Exec("DROP INDEX ?", clause.Column{Name: name}).Error
	})
}

func (m Migrator) CreateTable(values ...interface{}) (err error) {
	if err = m.Migrator.CreateTable(values...); err != nil {
		return
	}
	for _, value := range m.ReorderModels(values, false) {
		if err = m.RunWithValue(value, func(stmt *gorm.Statement) error {
			if stmt.Schema != nil {
				for _, field := range stmt.Schema.FieldsByDBName {
					if field.Comment != "" {
						if err := m.DB.Exec(
							"COMMENT ON COLUMN ?.? IS ?",
							m.CurrentTable(stmt), clause.Column{Name: field.DBName}, gorm.Expr(m.Migrator.Dialector.Explain("$1", field.Comment)),
						).Error; err != nil {
							return err
						}
					}
				}
			}
			return nil
		}); err != nil {
			return
		}
	}
	return
}

func (m Migrator) HasTable(value interface{}) bool {
	var count int64
	_ = m.RunWithValue(value, func(stmt *gorm.Statement) error {
		currentSchema, curTable := m.CurrentSchema(stmt, stmt.Table)
		return m.DB.Raw("SELECT count(*) FROM information_schema.tables WHERE table_schema = ? AND table_name = ? AND table_type = ?", currentSchema, curTable, "BASE TABLE").Scan(&count).Error
	})
	return count > 0
}

func (m Migrator) DropTable(values ...interface{}) error {
	values = m.ReorderModels(values, false)
	return m.DB.Connection(func(tx *gorm.DB) error {
		tx.Exec("SET FOREIGN_KEY_CHECKS = 0;")
		for i := len(values) - 1; i >= 0; i-- {
			if err := m.RunWithValue(values[i], func(stmt *gorm.Statement) error {
				return tx.Exec("DROP TABLE IF EXISTS ? CASCADE", clause.Table{Name: stmt.Table}).Error
			}); err != nil {
				return err
			}
		}
		return tx.Exec("SET FOREIGN_KEY_CHECKS = 1;").Error
	})
}

func (m Migrator) DropConstraint(value interface{}, name string) error {
	return m.RunWithValue(value, func(stmt *gorm.Statement) error {
		constraint, chk, table := m.GuessConstraintAndTable(stmt, name)
		if chk != nil {
			return m.DB.Exec("ALTER TABLE ? DROP CHECK ?", clause.Table{Name: stmt.Table}, clause.Column{Name: chk.Name}).Error
		}
		if constraint != nil {
			name = constraint.Name
		}

		return m.DB.Exec(
			"ALTER TABLE ? DROP FOREIGN KEY ?", clause.Table{Name: table}, clause.Column{Name: name},
		).Error
	})
}

func (m Migrator) AddColumn(value interface{}, field string) error {
	if err := m.Migrator.AddColumn(value, field); err != nil {
		return err
	}
	m.resetPreparedStmts()

	return m.RunWithValue(value, func(stmt *gorm.Statement) error {
		if stmt.Schema != nil {
			if field := stmt.Schema.LookUpField(field); field != nil {
				if field.Comment != "" {
					if err := m.DB.Exec(
						"COMMENT ON COLUMN ?.? IS ?",
						m.CurrentTable(stmt), clause.Column{Name: field.DBName}, gorm.Expr(m.Migrator.Dialector.Explain("$1", field.Comment)),
					).Error; err != nil {
						return err
					}
				}
			}
		}
		return nil
	})
}

func (m Migrator) HasColumn(value interface{}, field string) bool {
	var count int64
	_ = m.RunWithValue(value, func(stmt *gorm.Statement) error {
		name := field
		if stmt.Schema != nil {
			if field := stmt.Schema.LookUpField(field); field != nil {
				name = field.DBName
			}
		}

		currentSchema, curTable := m.CurrentSchema(stmt, stmt.Table)
		return m.DB.Raw(
			"SELECT count(*) FROM INFORMATION_SCHEMA.columns WHERE table_schema = ? AND table_name = ? AND column_name = ?",
			currentSchema, curTable, name,
		).Scan(&count).Error
	})

	return count > 0
}

func (m Migrator) MigrateColumn(value interface{}, field *schema.Field, columnType gorm.ColumnType) error {
	// skip primary field
	if !field.PrimaryKey {
		if err := m.Migrator.MigrateColumn(value, field, columnType); err != nil {
			return err
		}
	}

	return m.RunWithValue(value, func(stmt *gorm.Statement) error {
		var description string
		currentSchema, curTable := m.CurrentSchema(stmt, stmt.Table)
		values := []interface{}{currentSchema, curTable, field.DBName, stmt.Table, currentSchema}
		checkSQL := "SELECT description FROM pg_catalog.pg_description "
		checkSQL += "WHERE objsubid = (SELECT ordinal_position FROM information_schema.columns WHERE table_schema = ? AND table_name = ? AND column_name = ?) "
		checkSQL += "AND objoid = (SELECT oid FROM pg_catalog.pg_class WHERE relname = ? AND relnamespace = "
		checkSQL += "(SELECT oid FROM pg_catalog.pg_namespace WHERE nspname = ?))"
		m.DB.Raw(checkSQL, values...).Scan(&description)

		comment := strings.Trim(field.Comment, "'")
		comment = strings.Trim(comment, `"`)
		if field.Comment != "" && comment != description {
			if err := m.DB.Exec(
				"COMMENT ON COLUMN ?.? IS ?",
				m.CurrentTable(stmt), clause.Column{Name: field.DBName}, gorm.Expr(m.Migrator.Dialector.Explain("$1", field.Comment)),
			).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (m Migrator) HasConstraint(value interface{}, name string) bool {
	var count int64
	_ = m.RunWithValue(value, func(stmt *gorm.Statement) error {
		constraint, table := m.GuessConstraintInterfaceAndTable(stmt, name)
		currentSchema, curTable := m.CurrentSchema(stmt, table)
		if constraint != nil {
			name = constraint.GetName()
		}
		return m.DB.Raw(
			"SELECT count(*) FROM INFORMATION_SCHEMA.table_constraints WHERE table_schema = ? AND table_name = ? AND constraint_name = ?",
			currentSchema, curTable, name,
		).Scan(&count).Error
	})

	return count > 0
}
func (m Migrator) ColumnTypes(value interface{}) ([]gorm.ColumnType, error) {
	columnTypes := make([]gorm.ColumnType, 0)
	err := m.RunWithValue(value, func(stmt *gorm.Statement) error {
		var (
			currentDatabase, table = m.CurrentSchema(stmt, stmt.Table)
			columnTypeSQL          = "SELECT column_name, column_default, is_nullable = 'YES', data_type, character_maximum_length, column_type, column_key, extra, column_comment, numeric_precision, numeric_scale "
			rows, err              = m.DB.Session(&gorm.Session{}).Table(table).Limit(1).Rows()
		)
		if err != nil {
			return err
		}
		rawColumnTypes, err := rows.ColumnTypes()
		if err != nil {
			return err
		}
		if err := rows.Close(); err != nil {
			return err
		}
		columnTypeSQL += "FROM information_schema.columns WHERE table_schema = ? AND table_name = ? ORDER BY ORDINAL_POSITION"
		columns, rowErr := m.DB.Table(table).Raw(columnTypeSQL, currentDatabase, table).Rows()
		if rowErr != nil {
			return rowErr
		}
		defer columns.Close()
		for columns.Next() {
			var (
				column            migrator.ColumnType
				datetimePrecision sql.NullInt64
				extraValue        sql.NullString
				columnKey         sql.NullString
				values            = []interface{}{
					&column.NameValue, &column.DefaultValueValue, &column.NullableValue, &column.DataTypeValue, &column.LengthValue, &column.ColumnTypeValue, &columnKey, &extraValue, &column.CommentValue, &column.DecimalSizeValue, &column.ScaleValue,
				}
			)
			if scanErr := columns.Scan(values...); scanErr != nil {
				return scanErr
			}
			column.PrimaryKeyValue = sql.NullBool{Bool: false, Valid: true}
			column.UniqueValue = sql.NullBool{Bool: false, Valid: true}
			switch columnKey.String {
			case "PRI":
				column.PrimaryKeyValue = sql.NullBool{Bool: true, Valid: true}
			case "UNI":
				column.UniqueValue = sql.NullBool{Bool: true, Valid: true}
			}
			if strings.Contains(extraValue.String, "auto_increment") {
				column.AutoIncrementValue = sql.NullBool{Bool: true, Valid: true}
			}
			column.DefaultValueValue.String = strings.Trim(column.DefaultValueValue.String, "'")
			if datetimePrecision.Valid {
				column.DecimalSizeValue = datetimePrecision
			}
			for _, c := range rawColumnTypes {
				if c.Name() == column.NameValue.String {
					column.SQLColumnType = c
					break
				}
			}
			columnTypes = append(columnTypes, column)
		}
		return nil
	})
	return columnTypes, err
}

func (m Migrator) CreateSequence(tx *gorm.DB, stmt *gorm.Statement, field *schema.Field,
	serialDatabaseType string) (err error) {
	_, table := m.CurrentSchema(stmt, stmt.Table)
	tableName := table
	sequenceName := strings.Join([]string{tableName, field.DBName, "seq"}, "_")
	if err = tx.Exec(`CREATE SEQUENCE IF NOT EXISTS ? AS ?`, clause.Expr{SQL: sequenceName},
		clause.Expr{SQL: serialDatabaseType}).Error; err != nil {
		return err
	}
	if err := tx.Exec("ALTER TABLE ? ALTER COLUMN ? SET DEFAULT nextval('?')",
		clause.Expr{SQL: tableName}, clause.Expr{SQL: field.DBName}, clause.Expr{SQL: sequenceName}).Error; err != nil {
		return err
	}
	if err := tx.Exec("ALTER SEQUENCE ? OWNED BY ?.?",
		clause.Expr{SQL: sequenceName}, clause.Expr{SQL: tableName}, clause.Expr{SQL: field.DBName}).Error; err != nil {
		return err
	}
	return
}

func (m Migrator) UpdateSequence(tx *gorm.DB, stmt *gorm.Statement, field *schema.Field,
	serialDatabaseType string) (err error) {

	sequenceName, err := m.getColumnSequenceName(tx, stmt, field)
	if err != nil {
		return err
	}

	if err = tx.Exec(`ALTER SEQUENCE IF EXISTS ? AS ?`, clause.Expr{SQL: sequenceName}, clause.Expr{SQL: serialDatabaseType}).Error; err != nil {
		return err
	}

	if err := tx.Exec("ALTER TABLE ? ALTER COLUMN ? TYPE ?",
		m.CurrentTable(stmt), clause.Expr{SQL: field.DBName}, clause.Expr{SQL: serialDatabaseType}).Error; err != nil {
		return err
	}
	return
}

func (m Migrator) DeleteSequence(tx *gorm.DB, stmt *gorm.Statement, field *schema.Field,
	fileType clause.Expr) (err error) {

	sequenceName, err := m.getColumnSequenceName(tx, stmt, field)
	if err != nil {
		return err
	}

	if err := tx.Exec("ALTER TABLE ? ALTER COLUMN ? TYPE ?", m.CurrentTable(stmt), clause.Column{Name: field.DBName}, fileType).Error; err != nil {
		return err
	}

	if err := tx.Exec("ALTER TABLE ? ALTER COLUMN ? DROP DEFAULT",
		m.CurrentTable(stmt), clause.Expr{SQL: field.DBName}).Error; err != nil {
		return err
	}

	if err = tx.Exec(`DROP SEQUENCE IF EXISTS ?`, clause.Expr{SQL: sequenceName}).Error; err != nil {
		return err
	}

	return
}

func (m Migrator) getColumnSequenceName(tx *gorm.DB, stmt *gorm.Statement, field *schema.Field) (
	sequenceName string, err error) {
	_, table := m.CurrentSchema(stmt, stmt.Table)

	// DefaultValueValue is reset by ColumnTypes, search again.
	var columnDefault string
	err = tx.Raw(
		`SELECT column_default FROM information_schema.columns WHERE table_name = ? AND column_name = ?`,
		table, field.DBName).Scan(&columnDefault).Error

	if err != nil {
		return
	}

	sequenceName = strings.TrimSuffix(
		strings.TrimPrefix(columnDefault, `nextval('`),
		`'::regclass)`,
	)
	return
}

type Index struct {
	TableName  string `gorm:"column:TABLE_NAME"`
	ColumnName string `gorm:"column:COLUMN_NAME"`
	IndexName  string `gorm:"column:INDEX_NAME"`
	NonUnique  int32  `gorm:"column:NON_UNIQUE"`
}

func groupByIndexName(indexList []*Index) (map[string][]*Index, []string) {
	columnIndexMap := make(map[string][]*Index, len(indexList))
	indexNames := make([]string, 0, len(indexList))
	for _, idx := range indexList {
		if _, ok := columnIndexMap[idx.IndexName]; !ok {
			indexNames = append(indexNames, idx.IndexName)
		}
		columnIndexMap[idx.IndexName] = append(columnIndexMap[idx.IndexName], idx)
	}
	return columnIndexMap, indexNames
}

func (m Migrator) GetTypeAliases(databaseTypeName string) []string {
	return typeAliasMap[databaseTypeName]
}

// should reset prepared stmts when table changed
func (m Migrator) resetPreparedStmts() {
	if m.DB.PrepareStmt {
		if pdb, ok := m.DB.ConnPool.(*gorm.PreparedStmtDB); ok {
			pdb.Reset()
		}
	}
}

func (m Migrator) DropColumn(dst interface{}, field string) error {
	if err := m.Migrator.DropColumn(dst, field); err != nil {
		return err
	}

	m.resetPreparedStmts()
	return nil
}

func (m Migrator) RenameColumn(value interface{}, oldName, newName string) error {
	return m.RunWithValue(value, func(stmt *gorm.Statement) error {
		var field *schema.Field
		if stmt.Schema != nil {
			if f := stmt.Schema.LookUpField(oldName); f != nil {
				oldName = f.DBName
				field = f
			}

			if f := stmt.Schema.LookUpField(newName); f != nil {
				newName = f.DBName
				field = f
			}
		}

		if field != nil {
			return m.DB.Exec(
				"ALTER TABLE ? CHANGE ? ? ?",
				clause.Table{Name: stmt.Table}, clause.Column{Name: oldName},
				clause.Column{Name: newName}, m.FullDataTypeOf(field),
			).Error
		}

		return fmt.Errorf("failed to look up field with name: %s", newName)
	})
}

func (m Migrator) CurrentSchema(stmt *gorm.Statement, table string) (string, string) {
	if tables := strings.Split(table, `.`); len(tables) == 2 {
		return tables[0], tables[1]
	}
	m.DB = m.DB.Table(table)
	return m.CurrentDatabase(), table
}

// TableType table type return tableType,error
func (m Migrator) TableType(value interface{}) (tableType gorm.TableType, err error) {
	var table migrator.TableType
	err = m.RunWithValue(value, func(stmt *gorm.Statement) error {
		var (
			values = []interface{}{
				&table.SchemaValue, &table.NameValue, &table.TypeValue, &table.CommentValue,
			}
			currentDatabase, tableName = m.CurrentSchema(stmt, stmt.Table)
			tableTypeSQL               = "SELECT table_schema, table_name, table_type, table_comment FROM information_schema.tables WHERE table_schema = ? AND table_name = ?"
		)
		row := m.DB.Table(tableName).Raw(tableTypeSQL, currentDatabase, tableName).Row()
		if scanErr := row.Scan(values...); scanErr != nil {
			return scanErr
		}
		return nil
	})
	return table, err
}
