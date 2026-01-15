package dameng

import (
	"github.com/thoas/go-funk"
	"gorm.io/gorm"
	"gorm.io/gorm/callbacks"
	"gorm.io/gorm/clause"
	gormSchema "gorm.io/gorm/schema"
)

func Create(db *gorm.DB) {
	boundVars := make(map[string]int)
	if db.Statement == nil || db.Statement.Schema == nil {
		return
	}
	for k, v := range db.Statement.Clauses {
		if v.Name == "INSERT" && v.Expression.(clause.Insert).Modifier == "IGNORE" {
			delete(db.Statement.Clauses, k)
		}
	}
	if !db.Statement.Unscoped {
		for _, c := range db.Statement.Schema.CreateClauses {
			db.Statement.AddClause(c)
		}
	}
	if db.Statement.SQL.String() == "" {
		hasDefaultValues := len(db.Statement.Schema.FieldsWithDefaultDBValue) > 0
		values := callbacks.ConvertToCreateValues(db.Statement)
		onConflict, hasConflict := db.Statement.Clauses["ON CONFLICT"].Expression.(clause.OnConflict)
		if hasConflict {
			if len(onConflict.Columns) != len(db.Statement.Schema.PrimaryFields) {
				for _, v := range db.Statement.Schema.PrimaryFields {
					onConflict.Columns = append(onConflict.Columns, clause.Column{Name: v.DBName})
				}
			}
			hasConflict = contains(values.Columns, onConflict.Columns)
		}
		if hasConflict {
			damengCreateMerge(db, onConflict, values)
		} else {
			createInsert(db, hasDefaultValues, values, boundVars)
		}
		run(db, hasDefaultValues, values, boundVars)
	}
}
func contains(all, sub []clause.Column) bool {
	kv := make(map[string]struct{}, len(all))
	for _, v := range all {
		kv[v.Name] = struct{}{}
	}
	for _, v := range sub {
		if _, ok := kv[v.Name]; !ok {
			return false
		}
	}
	return true
}
func damengCreateMerge(db *gorm.DB, onConflict clause.OnConflict, values clause.Values) {
	_, _ = db.Statement.WriteString("MERGE INTO ")
	db.Statement.WriteQuoted(db.Statement.Table)
	_, _ = db.Statement.WriteString(" USING (")
	for i, vs := range values.Values {
		if i > 0 {
			_, _ = db.Statement.WriteString(" UNION ")
		}
		_, _ = db.Statement.WriteString("SELECT ")
		for j, v := range vs {
			if j > 0 {
				_ = db.Statement.WriteByte(',')
			}
			db.Statement.AddVar(db.Statement, v)
			_, _ = db.Statement.WriteString(" AS ")
			db.Statement.WriteQuoted(values.Columns[j].Name)
		}
		_, _ = db.Statement.WriteString(" FROM ")
		_, _ = db.Statement.WriteString(db.Dialector.(*Dialector).DummyTableName())
	}
	_, _ = db.Statement.WriteString(") AS excluded ON (")
	var where clause.Where
	colkv := make(map[string]struct{}, len(onConflict.Columns))
	for _, field := range onConflict.Columns {
		where.Exprs = append(where.Exprs, clause.Eq{
			Column: clause.Column{Table: db.Statement.Table, Name: field.Name},
			Value:  clause.Column{Table: "excluded", Name: field.Name},
		})
		colkv[field.Name] = struct{}{}
	}
	where.Build(db.Statement)
	_, _ = db.Statement.WriteString(")")
	if len(onConflict.DoUpdates) > 0 {
		_, _ = db.Statement.WriteString(" WHEN MATCHED THEN UPDATE SET ")
		var newUpdates clause.Set
		for _, v := range onConflict.DoUpdates {
			if _, ok := colkv[v.Column.Name]; ok {
				continue
			}
			switch v.Value.(type) {
			case clause.Expr:
				o := v.Value.(clause.Expr)
				o.SQL = db.Statement.Table + "." + o.SQL
				v.Value = o
			}
			v.Column.Table = db.Statement.Table
			newUpdates = append(newUpdates, v)
		}
		newUpdates.Build(db.Statement)
	}
	_, _ = db.Statement.WriteString(" WHEN NOT MATCHED THEN INSERT (")
	written := false
	if db.Statement.Schema.PrioritizedPrimaryField != nil {
		ai, ok := db.Statement.Schema.PrioritizedPrimaryField.TagSettings["AUTOINCREMENT"]
		db.Statement.Schema.PrioritizedPrimaryField.AutoIncrement = false
		if ok && ai == "true" {
			db.Statement.Schema.PrioritizedPrimaryField.AutoIncrement = true
		}
	}
	for _, column := range values.Columns {
		if db.Statement.Schema.PrioritizedPrimaryField == nil || !db.Statement.Schema.PrioritizedPrimaryField.AutoIncrement || db.Statement.Schema.PrioritizedPrimaryField.DBName != column.Name {
			if written {
				_ = db.Statement.WriteByte(',')
			}
			written = true
			db.Statement.WriteQuoted(column.Name)
		}
	}
	_, _ = db.Statement.WriteString(") VALUES (")
	written = false
	for _, column := range values.Columns {
		if db.Statement.Schema.PrioritizedPrimaryField == nil || !db.Statement.Schema.PrioritizedPrimaryField.AutoIncrement || db.Statement.Schema.PrioritizedPrimaryField.DBName != column.Name {
			if written {
				_ = db.Statement.WriteByte(',')
			}
			written = true
			db.Statement.WriteQuoted(clause.Column{Table: "excluded", Name: column.Name})
		}
	}
	_, _ = db.Statement.WriteString(");")
}
func createInsert(db *gorm.DB, hasDefaultValues bool, values clause.Values, boundVars map[string]int) {
	db.Statement.AddClauseIfNotExists(clause.Insert{Table: clause.Table{Name: db.Statement.Table}})
	db.Statement.AddClause(values)
	if hasDefaultValues {
		db.Statement.AddClauseIfNotExists(clause.Returning{
			Columns: funk.Map(db.Statement.Schema.FieldsWithDefaultDBValue, func(field *gormSchema.Field) clause.Column {
				return clause.Column{Name: field.DBName}
			}).([]clause.Column),
		})
	}
	db.Statement.Build("INSERT", "VALUES")
}
func run(db *gorm.DB, hasDefaultValues bool, values clause.Values, boundVars map[string]int) {
	if db.DryRun || db.Error != nil {
		return
	}
	result, err := db.Statement.ConnPool.ExecContext(db.Statement.Context, db.Statement.SQL.String(), db.Statement.Vars...)
	if db.AddError(err) == nil {
		db.RowsAffected, _ = result.RowsAffected()
	}
}
