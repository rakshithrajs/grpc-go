package storage

import (
	"fmt"
	"strings"
)

// UpdateField describes a single column/value pair used to build a dynamic UPDATE statement.
type UpdateField struct {
	Column string
	Value  any
}

// BuildUpdateSQL constructs a dynamic UPDATE statement for the given table, fields, and WHERE columns.
func BuildUpdateSQL(table string, fields []UpdateField, whereColumns []string) (string, []any) {
	args := make([]any, 0, len(whereColumns)+len(fields))
	for range whereColumns {
		args = append(args, nil)
	}

	setParts := make([]string, 0, len(fields)+1)
	for _, f := range fields {
		args = append(args, f.Value)
		setParts = append(setParts, fmt.Sprintf(`"%s" = $%d`, f.Column, len(args)))
	}

	setParts = append(setParts, `"updatedAtUTC" = NOW()`)

	whereParts := make([]string, len(whereColumns))
	for i, col := range whereColumns {
		whereParts[i] = fmt.Sprintf(`"%s" = $%d`, col, i+1)
	}

	query := fmt.Sprintf(`UPDATE %s SET %s WHERE %s`, table, strings.Join(setParts, ", "), strings.Join(whereParts, " AND "))

	return query, args
}
