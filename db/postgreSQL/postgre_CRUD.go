package postgreSQL

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"GO_Music/db"
)

type PostgresRepository[T any, ID comparable] struct {
	db        *sql.DB
	tx        *sql.Tx
	tableName string
	idColumn  string
}

func NewPostgresRepository[T any, ID comparable](db *sql.DB, tableName string, idColumn string) *PostgresRepository[T, ID] {
	return &PostgresRepository[T, ID]{
		db:        db,
		tableName: tableName,
		idColumn:  idColumn,
	}
}

func (r *PostgresRepository[T, ID]) WithTx(tx *sql.Tx) db.Repository[T, ID] {
	return &PostgresRepository[T, ID]{
		db:        r.db,
		tx:        tx,
		tableName: r.tableName,
		idColumn:  r.idColumn,
	}
}

func (r *PostgresRepository[T, ID]) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	if r.tx != nil {
		return r.tx.ExecContext(ctx, query, args...)
	}
	return r.db.ExecContext(ctx, query, args...)
}

func (r *PostgresRepository[T, ID]) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	if r.tx != nil {
		return r.tx.QueryRowContext(ctx, query, args...)
	}
	return r.db.QueryRowContext(ctx, query, args...)
}

func (r *PostgresRepository[T, ID]) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	if r.tx != nil {
		return r.tx.QueryContext(ctx, query, args...)
	}
	return r.db.QueryContext(ctx, query, args...)
}

// Create конвертирует struct в map через reflect и вставляет запись
func (r *PostgresRepository[T, ID]) Create(ctx context.Context, entity *T) error {
	m, err := db.StructToMap(entity)
	if err != nil {
		return err
	}

	// Исключаем ID-колонку из вставки
	delete(m, r.idColumn)

	columns := make([]string, 0, len(m))
	placeholders := make([]string, 0, len(m))
	values := make([]interface{}, 0, len(m))

	i := 1
	for col, val := range m {
		columns = append(columns, col)
		placeholders = append(placeholders, fmt.Sprintf("$%d", i))
		values = append(values, val)
		i++
	}

	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s) RETURNING %s",
		r.tableName,
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", "),
		r.idColumn,
	)

	// Получаем сгенерированный ID
	var id ID
	err = r.QueryRowContext(ctx, query, values...).Scan(&id)
	if err != nil {
		return err
	}

	// Устанавливаем полученный ID обратно в структуру
	if setter, ok := any(entity).(interface{ SetID(ID) }); ok {
		setter.SetID(id)
	}

	return nil
}

// Update обновляет запись по id
func (r *PostgresRepository[T, ID]) Update(ctx context.Context, entity *T) error {
	m, err := db.StructToMap(entity)
	if err != nil {
		return err
	}

	idVal, ok := m[r.idColumn]
	if !ok {
		return fmt.Errorf("entity must have field %s", r.idColumn)
	}
	delete(m, r.idColumn)

	setParts := make([]string, 0, len(m))
	values := make([]interface{}, 0, len(m)+1)

	i := 1
	for col, val := range m {
		setParts = append(setParts, fmt.Sprintf("%s = $%d", col, i))
		values = append(values, val)
		i++
	}
	values = append(values, idVal)

	query := fmt.Sprintf(
		"UPDATE %s SET %s WHERE %s = $%d",
		r.tableName,
		strings.Join(setParts, ", "),
		r.idColumn,
		i,
	)
	_, err = r.ExecContext(ctx, query, values...)
	return err
}

func (r *PostgresRepository[T, ID]) Delete(ctx context.Context, id ID) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE %s = $1", r.tableName, r.idColumn)
	_, err := r.ExecContext(ctx, query, id)
	return err
}

func (r *PostgresRepository[T, ID]) GetByID(ctx context.Context, id ID) (*T, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s = $1", r.tableName, r.idColumn)
	row := r.QueryRowContext(ctx, query, id)

	cols, err := r.getColumns(ctx)
	if err != nil {
		return nil, err
	}

	// Создаем слайс для сканирования, учитывая nullable поля
	values := make([]interface{}, len(cols))
	valuePtrs := make([]interface{}, len(cols))

	// Для nullable строк используем sql.NullString
	nullStringMap := map[string]*sql.NullString{}
	for i, col := range cols {
		if strings.HasSuffix(col, "_string") { // или другой способ идентификации
			ns := new(sql.NullString)
			values[i] = ns
			nullStringMap[col] = ns
			valuePtrs[i] = ns
		} else {
			var v interface{}
			values[i] = &v
			valuePtrs[i] = &v
		}
	}

	err = row.Scan(valuePtrs...)
	if err != nil {
		return nil, err
	}

	m := make(map[string]interface{})
	for i, col := range cols {
		if ns, ok := nullStringMap[col]; ok {
			if ns.Valid {
				// Создаем копию строки, чтобы избежать проблем с указателями
				str := ns.String
				m[col] = &str
			} else {
				m[col] = nil
			}
		} else {
			valPtr := values[i].(*interface{})
			val := *valPtr
			if b, ok := val.([]byte); ok {
				m[col] = string(b)
			} else {
				m[col] = val
			}
		}
	}

	var entity T
	if err := db.MapToStruct(m, &entity); err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *PostgresRepository[T, ID]) GetByIDs(ctx context.Context, ids []ID) ([]*T, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	params := make([]string, len(ids))
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		params[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}

	query := fmt.Sprintf("SELECT * FROM %s WHERE %s IN (%s)", r.tableName, r.idColumn, strings.Join(params, ", "))
	rows, err := r.QueryContext(ctx, query, args...)
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
