package postgreSQL

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"GO_Music/db"
)

func (r *PostgresRepository[T, ID]) List(ctx context.Context, filter db.Filter) ([]*T, error) {
	queryBuilder := strings.Builder{}
	queryBuilder.WriteString(fmt.Sprintf("SELECT * FROM %s", r.tableName))

	args := []interface{}{}
	argPos := 1

	if len(filter.Conditions) > 0 {
		conds := make([]string, 0, len(filter.Conditions))
		for _, cond := range filter.Conditions {
			conds = append(conds, fmt.Sprintf("%s %s $%d", cond.Field, cond.Operator, argPos))
			args = append(args, cond.Value)
			argPos++
		}
		queryBuilder.WriteString(" WHERE " + strings.Join(conds, " AND "))
	}

	// ORDER BY
	if filter.OrderBy != "" {
		queryBuilder.WriteString(" ORDER BY " + filter.OrderBy)
	}

	// LIMIT OFFSET
	if filter.Limit > 0 {
		queryBuilder.WriteString(" LIMIT " + strconv.Itoa(filter.Limit))
	}
	if filter.Offset > 0 {
		queryBuilder.WriteString(" OFFSET " + strconv.Itoa(filter.Offset))
	}

	query := queryBuilder.String()

	rows, err := r.queryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	var results []*T
	for rows.Next() {
		values := make([]interface{}, len(cols))
		valuePtrs := make([]interface{}, len(cols))
		for i := range cols {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, err
		}

		m := make(map[string]interface{})
		for i, col := range cols {
			val := values[i]
			if b, ok := val.([]byte); ok {
				m[col] = string(b)
			} else {
				m[col] = val
			}
		}

		var entity T
		if err := db.MapToStruct(m, &entity); err != nil {
			return nil, err
		}
		results = append(results, &entity)
	}

	return results, nil
}

func (r *PostgresRepository[T, ID]) Count(ctx context.Context, filter db.Filter) (int, error) {
	queryBuilder := strings.Builder{}
	queryBuilder.WriteString(fmt.Sprintf("SELECT COUNT(*) FROM %s", r.tableName))

	args := []interface{}{}
	argPos := 1

	if len(filter.Conditions) > 0 {
		conds := make([]string, 0, len(filter.Conditions))
		for _, cond := range filter.Conditions {
			conds = append(conds, fmt.Sprintf("%s %s $%d", cond.Field, cond.Operator, argPos))
			args = append(args, cond.Value)
			argPos++
		}
		queryBuilder.WriteString(" WHERE " + strings.Join(conds, " AND "))
	}

	query := queryBuilder.String()

	var count int
	err := r.queryRowContext(ctx, query, args...).Scan(&count)
	return count, err
}

func (r *PostgresRepository[T, ID]) Exists(ctx context.Context, id ID) (bool, error) {
	query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE %s = $1)", r.tableName, r.idColumn)
	var exists bool
	err := r.queryRowContext(ctx, query, id).Scan(&exists)
	return exists, err
}

// getColumns возвращает список колонок таблицы
func (r *PostgresRepository[T, ID]) getColumns(ctx context.Context) ([]string, error) {
	query := fmt.Sprintf("SELECT column_name FROM information_schema.columns WHERE table_name = '%s' ORDER BY ordinal_position", r.tableName)
	rows, err := r.queryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var columns []string
	for rows.Next() {
		var col string
		if err := rows.Scan(&col); err != nil {
			return nil, err
		}
		columns = append(columns, col)
	}
	return columns, nil
}
