package dbutil

import (
	"database/sql"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
)

func get_columns_mysql(db *sql.DB, st string) ([]string, error) {
	slice := strings.Split(st, ".")
	if len(slice) != 2 {
		return nil, fmt.Errorf("参数错误 schema table")
	}
	query := `
	select column_name
	from information_schema.columns
	where table_schema = ? and table_name = ?
	order by ordinal_position;
	`
	rows, err := db.Query(query, slice[0], slice[1])
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	columns := []string{}
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		columns = append(columns, name)
	}
	return columns, nil
}

type MySQLRepoImpl struct {
	db *sql.DB
}

// NewMySQLRepo creates a new MySQLRepoImpl instance.
//
// Parameters:
//   - db: database connection
//
// Returns:
//   - *MySQLRepoImpl: MySQLRepoImpl instance
func NewMySQLRepo(db *sql.DB) *MySQLRepoImpl {
	return &MySQLRepoImpl{db: db}
}

// Create inserts a new record into the specified table (MySQL).
//
// Parameters:
//   - st: schema and table, format like "schema.table"
//   - d: data to be inserted
//
// Returns:
//   - error: error information
func (r *MySQLRepoImpl) Create(st string, d map[string]interface{}) error {
	columns, err := get_columns_mysql(r.db, st)
	if err != nil {
		return err
	}

	var placeholders []string
	var values []interface{}
	for _, column := range columns {
		if val, ok := d[column]; ok {
			placeholders = append(placeholders, "?")
			values = append(values, val)
		}
	}

	q := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", st, strings.Join(columns, ", "), strings.Join(placeholders, ", "))
	stmt, err := r.db.Prepare(q)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(values...)
	return err
}

// Get retrieves records from the specified table based on conditions.
//
// Parameters:
//   - st: schema and table, format like "schema.table"
//   - c: columns to retrieve, e.g., ["id", "name"]
//   - f: filter conditions, e.g., [["equal", "name", "John Doe"], ["in", "id", "1a", "1b"]]
//   - l: additional clauses, e.g., "order by id desc limit 20 offset 0"
//
// Returns:
//   - []map[string]interface{}: retrieved records
//   - error: error information
func (r *MySQLRepoImpl) Get(st string, c []string, f [][]string, l string) ([]map[string]interface{}, error) {
	if len(c) == 0 {
		var err error
		c, err = get_columns_mysql(r.db, st)
		if err != nil {
			log.Println(err.Error())
			return nil, err
		}
	}
	q := fmt.Sprintf("SELECT %s FROM %s", strings.Join(c, ", "), st)

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
			whereClauses = append(whereClauses, fmt.Sprintf("POSITION(? IN %s) > 0", field))
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
		case "json-array-contains":
			whereClauses = append(whereClauses, fmt.Sprintf("JSON_CONTAINS(%s, '\"%s\"')", field, condition[2]))
		case "json-object-contains":
			whereClauses = append(whereClauses, fmt.Sprintf("JSON_CONTAINS(%s, ?, '$')", field))
			params = append(params, fmt.Sprintf(`{"%s": "%s"}`, condition[2], condition[3]))
		}
	}

	if len(whereClauses) > 0 {
		q += " WHERE " + strings.Join(whereClauses, " AND ")
	}

	if l != "" {
		q += " " + l
	}

	log.Println(q)
	stmt, err := r.db.Prepare(q)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	log.Println(params)
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

// Update modifies records in the specified table based on conditions (MySQL).
//
// Parameters:
//   - st: schema and table, format like "schema.table"
//   - d: data to be updated
//   - w: WHERE condition, e.g., "id='1a'"
//
// Returns:
//   - error: error information
func (r *MySQLRepoImpl) Update(st string, d map[string]interface{}, w string) error {
	columns, err := get_columns_mysql(r.db, st)
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

// Remove deletes records from the specified table based on conditions (MySQL).
//
// Parameters:
//   - st: schema and table, format like "schema.table"
//   - w: WHERE condition, e.g., "id='1a'"
//
// Returns:
//   - error: error information
func (r *MySQLRepoImpl) Remove(st string, w string) error {
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
