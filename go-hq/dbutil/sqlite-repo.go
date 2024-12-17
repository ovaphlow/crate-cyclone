package dbutil

import (
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// get_columns_sqlite retrieves the column names of a given SQLite table.
// Parameters:
// - db: The database connection.
// - sat: The name of the table.
// Returns:
// - A slice of column names.
// - An error if the query fails.
func get_columns_sqlite(db *sql.DB, sat string) ([]string, error) {
	rows, err := db.Query("PRAGMA table_info(" + sat + ")")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var columns []string
	for rows.Next() {
		var cid int
		var name string
		var dtype string
		var notnull int
		var dflt sql.NullString // 修改这里
		var pk int
		err = rows.Scan(&cid, &name, &dtype, &notnull, &dflt, &pk)
		if err != nil {
			return nil, err
		}
		columns = append(columns, name)
	}
	return columns, nil
}

type SQLiteRepoImpl struct {
	db *sql.DB
}

// NewSQLiteRepo creates a new SQLiteRepoImpl instance.
// Parameters:
// - db: The database connection.
// Returns:
// - A pointer to the new SQLiteRepoImpl instance.
func NewSQLiteRepo(db *sql.DB) *SQLiteRepoImpl {
	return &SQLiteRepoImpl{db: db}
}

// Create inserts a new record into the specified table.
// Parameters:
// - st: The name of the table.
// - d: A map of column names to values.
// Returns:
// - An error if the operation fails.
func (r *SQLiteRepoImpl) Create(st string, d map[string]interface{}) error {
	columns, err := get_columns_sqlite(r.db, st)
	if err != nil {
		return err
	}

	var columnStr string
	var placeholders string
	var values []interface{}
	for _, column := range columns {
		if val, ok := d[column]; ok {
			columnStr += column + ","
			placeholders += "?,"
			values = append(values, val)
		}
	}
	columnStr = columnStr[:len(columnStr)-1]
	placeholders = placeholders[:len(placeholders)-1]

	stmt, err := r.db.Prepare("INSERT INTO " + st + " (" + columnStr + ") VALUES (" + placeholders + ")")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(values...)
	return err
}

// Get retrieves records from the specified table.
// Parameters:
// - st: The name of the table.
// - c: A slice of column names to retrieve.
// - f: A slice of filter conditions.
// - l: Additional SQL clauses (e.g., ORDER BY).
// Returns:
// - A slice of maps representing the retrieved records.
// - An error if the operation fails.
func (r *SQLiteRepoImpl) Get(st string, c []string, f [][]string, l string) ([]map[string]interface{}, error) {
	if len(c) == 0 {
		var err error
		c, err = get_columns_sqlite(r.db, st)
		if err != nil {
			return nil, err
		}
	}
	q := fmt.Sprintf("SELECT %s FROM %s", strings.Join(c, ","), st)

	var whereClauses []string
	var params []interface{}

	for _, condition := range f {
		if len(condition) < 3 {
			continue
		}
		field := condition[1]
		operator := condition[0]
		switch operator {
		case "equal":
			whereClauses = append(whereClauses, fmt.Sprintf("%s = ?", field))
			params = append(params, condition[2])
		case "not-equal":
			whereClauses = append(whereClauses, fmt.Sprintf("%s != ?", field))
			params = append(params, condition[2])
		case "in":
			placeholders := strings.Repeat("?,", len(condition)-2)
			placeholders = placeholders[:len(placeholders)-1]
			whereClauses = append(whereClauses, fmt.Sprintf("%s IN (%s)", field, placeholders))
			for _, v := range condition[2:] {
				params = append(params, v)
			}
		case "not-in":
			placeholders := strings.Repeat("?,", len(condition)-2)
			placeholders = placeholders[:len(placeholders)-1]
			whereClauses = append(whereClauses, fmt.Sprintf("%s NOT IN (%s)", field, placeholders))
			for _, v := range condition[2:] {
				params = append(params, v)
			}
		case "like":
			whereClauses = append(whereClauses, fmt.Sprintf("%s LIKE ?", field))
			params = append(params, condition[2])
		case "greater":
			whereClauses = append(whereClauses, fmt.Sprintf("%s > ?", field))
			params = append(params, condition[2])
		case "greater-equal":
			whereClauses = append(whereClauses, fmt.Sprintf("%s >= ?", field))
			params = append(params, condition[2])
		case "less":
			whereClauses = append(whereClauses, fmt.Sprintf("%s < ?", field))
			params = append(params, condition[2])
		case "less-equal":
			whereClauses = append(whereClauses, fmt.Sprintf("%s <= ?", field))
			params = append(params, condition[2])
		}
	}

	if len(whereClauses) > 0 {
		q += " WHERE " + strings.Join(whereClauses, " AND ")
	}

	if l != "" {
		q += " " + l
	}

	stmt, err := r.db.Prepare(q)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	var result []map[string]interface{}
	values := make([]interface{}, len(columns))
	valuePtrs := make([]interface{}, len(columns))
	for rows.Next() {
		for i := range columns {
			valuePtrs[i] = &values[i]
		}
		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, err
		}
		m := make(map[string]interface{})
		for i, col := range columns {
			val := values[i]
			if val == nil {
				m[col] = nil
			} else {
				switch v := val.(type) {
				case []byte:
					m[col] = string(v)
				case int, int8, int16, int32, int64:
					m[col] = strconv.FormatInt(reflect.ValueOf(v).Int(), 10)
				case uint, uint8, uint16, uint32, uint64:
					m[col] = strconv.FormatUint(reflect.ValueOf(v).Uint(), 10)
				default:
					m[col] = v
				}
			}
		}
		result = append(result, m)
	}

	return result, nil
}

// Update modifies existing records in the specified table.
// Parameters:
// - st: The name of the table.
// - d: A map of column names to new values.
// - w: The WHERE clause to specify which records to update.
// - deprecated: A boolean flag for deprecated usage.
// Returns:
// - An error if the operation fails.
func (r *SQLiteRepoImpl) Update(st string, d map[string]interface{}, w string) error {
	columns, err := get_columns_sqlite(r.db, st)
	if err != nil {
		return err
	}

	q := fmt.Sprintf("UPDATE %s SET ", st)
	var assignments []string
	var values []interface{}
	for _, column := range columns {
		if val, ok := d[column]; ok {
			assignments = append(assignments, fmt.Sprintf("%s = ?", column))
			values = append(values, val)
		}
	}
	q += strings.Join(assignments, ", ")
	q += " WHERE " + w

	stmt, err := r.db.Prepare(q)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(values...)
	return err
}

// Remove deletes records from the specified table.
// Parameters:
// - st: The name of the table.
// - w: The WHERE clause to specify which records to delete.
// Returns:
// - An error if the operation fails.
func (r *SQLiteRepoImpl) Remove(st string, w string) error {
	q := fmt.Sprintf("DELETE FROM %s WHERE %s", st, w)
	stmt, err := r.db.Prepare(q)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec()
	if err != nil {
		return err
	}
	return nil
}
