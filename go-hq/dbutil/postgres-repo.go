package dbutil

import (
	"database/sql"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
)

// get_columns_postgres retrieves column names for a given schema and table.
// Parameters:
// - db: database connection
// - sat: schema and table in "schema.table" format
// Returns:
// - []string: list of column names
// - error: error information
func get_columns_postgres(db *sql.DB, sat string) ([]string, error) {
	st := strings.Split(sat, ".")
	if len(st) != 2 {
		return []string{"*"}, nil
	}
	columns := []string{}
	stmt, err := db.Prepare(`
	SELECT column_name FROM information_schema.columns
	WHERE table_schema = $1 AND table_name = $2
	ORDER BY ordinal_position ASC
	`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(st[0], st[1])
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var column string
		err := rows.Scan(&column)
		if err != nil {
			return nil, err
		}
		columns = append(columns, column)
	}
	return columns, nil
}

type PostgresRepoImpl struct {
	db *sql.DB
}

// NewPostgresRepo creates a new PostgresRepoImpl instance.
// Parameters:
// - db: database connection
// Returns:
// - *PostgresRepoImpl: PostgresRepoImpl instance
func NewPostgresRepo(db *sql.DB) *PostgresRepoImpl {
	return &PostgresRepoImpl{db: db}
}

// Create inserts a new record into the specified table.
// Parameters:
// - st: schema and table in "schema.table" format
// - d: data to insert
// Returns:
// - error: error information
func (r *PostgresRepoImpl) Create(st string, d map[string]interface{}) error {
	columns, err := get_columns_postgres(r.db, st)
	if err != nil {
		return err
	}

	var values []string
	for _, column := range columns {
		if _, ok := d[column]; ok {
			values = append(values, fmt.Sprintf("%v", d[column]))
		}
	}

	q := fmt.Sprintf("insert into %s (%s) values (", st, strings.Join(columns, ", "))
	if len(values) == 0 {
		return nil
	}
	for i := 0; i < len(values); i++ {
		q += "$" + strconv.Itoa(i+1)
		if i < len(values)-1 {
			q += ","
		}
	}
	q += ")"
	p := make([]interface{}, len(values))
	for i, v := range values {
		p[i] = v
	}

	stmt, err := r.db.Prepare(q)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(p...)
	return err
}

// Get retrieves records from the specified table based on conditions.
// Parameters:
// - st: schema and table in "schema.table" format
// - c: columns to retrieve, e.g., ["id", "name"]
// - f: filter conditions, e.g., [["equal", "name", "John Doe"], ["in", "id", "1a", "1b"]]
// - l: additional clauses, e.g., "order by id desc limit 20 offset 0"
// Returns:
// - []map[string]interface{}: retrieved records
// - error: error information
func (r *PostgresRepoImpl) Get(st string, c []string, f [][]string, l string) ([]map[string]interface{}, error) {
	if len(c) == 0 {
		var err error
		c, err = get_columns_postgres(r.db, st)
		if err != nil {
			return nil, err
		}
	}
	q := fmt.Sprintf("select %s from %s", strings.Join(c, ", "), st)

	var whereClauses []string
	var params []interface{}
	paramIndex := 1

	for _, condition := range f {
		if len(condition) < 3 {
			continue
		}
		field := condition[1]
		operator := condition[0]
		switch operator {
		case "equal":
			whereClauses = append(whereClauses, fmt.Sprintf("%s = $%d", field, paramIndex))
			params = append(params, condition[2])
			paramIndex++
		case "not-equal":
			whereClauses = append(whereClauses, fmt.Sprintf("%s != $%d", field, paramIndex))
			params = append(params, condition[2])
			paramIndex++
		case "in":
			placeholders := make([]string, len(condition)-2)
			for i := range placeholders {
				placeholders[i] = fmt.Sprintf("$%d", paramIndex)
				params = append(params, condition[i+2])
				paramIndex++
			}
			whereClauses = append(whereClauses, fmt.Sprintf("%s in (%s)", field, strings.Join(placeholders, ", ")))
		case "not-in":
			placeholders := make([]string, len(condition)-2)
			for i := range placeholders {
				placeholders[i] = fmt.Sprintf("$%d", paramIndex)
				params = append(params, condition[i+2])
				paramIndex++
			}
			whereClauses = append(whereClauses, fmt.Sprintf("%s not in (%s)", field, strings.Join(placeholders, ", ")))
		case "greater":
			whereClauses = append(whereClauses, fmt.Sprintf("%s > $%d", field, paramIndex))
			params = append(params, condition[2])
			paramIndex++
		case "greater-equal":
			whereClauses = append(whereClauses, fmt.Sprintf("%s >= $%d", field, paramIndex))
			params = append(params, condition[2])
			paramIndex++
		case "less":
			whereClauses = append(whereClauses, fmt.Sprintf("%s < $%d", field, paramIndex))
			params = append(params, condition[2])
			paramIndex++
		case "less-equal":
			whereClauses = append(whereClauses, fmt.Sprintf("%s <= $%d", field, paramIndex))
			params = append(params, condition[2])
			paramIndex++
		case "jsonb-array-contains":
			whereClauses = append(whereClauses, fmt.Sprintf("%s @> '["+"%s"+"]'::jsonb", field, condition[2]))
			params = append(params, condition[2])
			paramIndex++
		case "jsonb-object-contains":
			whereClauses = append(whereClauses, fmt.Sprintf("%s @> $%d::jsonb", field, paramIndex))
			params = append(params, fmt.Sprintf(`{"%s": "%s"}`, condition[1], condition[2]))
			paramIndex++
		}
	}

	if len(whereClauses) > 0 {
		q += " where " + strings.Join(whereClauses, " AND ")
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

// Update modifies records in the specified table based on conditions.
// Parameters:
// - st: schema and table in "schema.table" format
// - d: data to update
// - w: WHERE condition, e.g., "id='1a'"
// Returns:
// - error: error information
func (r *PostgresRepoImpl) Update(st string, d map[string]interface{}, w string) error {
	columns, err := get_columns_postgres(r.db, st)
	if err != nil {
		return err
	}

	q := fmt.Sprintf("update %s set ", st)
	var values []string
	var p []interface{}
	for _, v := range columns {
		if _, ok := d[v]; ok {
			values = append(values, fmt.Sprintf("%s = $%d", v, len(values)+1))
			p = append(p, d[v])
		}
	}
	q += strings.Join(values, ", ")
	q += " where " + w

	log.Println(q)
	stmt, err := r.db.Prepare(q)
	if err != nil {
		return err
	}
	defer stmt.Close()

	log.Println(p)
	_, err = stmt.Exec(p...)
	if err != nil {
		return err
	}

	return nil
}

// Remove deletes records from the specified table based on conditions.
// Parameters:
// - st: schema and table in "schema.table" format
// - w: WHERE condition, e.g., "id='1a'"
// Returns:
// - error: error information
func (r *PostgresRepoImpl) Remove(st string, w string) error {
	q := fmt.Sprintf("delete from %s where %s", st, w)
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
