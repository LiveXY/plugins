package clauses

import "gorm.io/gorm/clause"

// IN Whether a value is within a set of values
type IN struct {
	Column any
	Values []any
}

func (in IN) Build(builder clause.Builder) {
	builder.WriteQuoted(in.Column)
	switch len(in.Values) {
	case 0:
		_, _ = builder.WriteString(" IN (NULL)")
	case 1:
		if _, ok := in.Column.([]clause.Column); ok {
			_, _ = builder.WriteString(" = (")
			builder.AddVar(builder, in.Values...)
			_, _ = builder.WriteString(")")
		} else {
			_, _ = builder.WriteString(" = ")
			builder.AddVar(builder, in.Values...)
		}

	default:
		_, _ = builder.WriteString(" IN (")
		builder.AddVar(builder, in.Values...)
		_ = builder.WriteByte(')')
	}
}
