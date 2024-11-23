package dbutil

import (
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// get_columns 获取指定 schema 和表的列名。
// 参数：
// - db：数据库连接
// - sat：schema 和表，格式如 "schema.table"
// 返回值：
// - []string：列名列表
// - error：错误信息
func get_columns(db *sql.DB, sat string) ([]string, error) {
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

type SharedRepo interface {
	// Create 插入一条新记录到指定的表。
	// 参数：
	// - st：schema 和表，格式如 "schema.table"
	// - d：要插入的数���
	// 返回值：
	// - error：错误信息
	Create(st string, d map[string]interface{}) error
	// Get 根据条件从指定的表中检索记录。
	// 参数：
	// - st：schema 和表，格式如 "schema.table"
	// - c：要检索的列，例如 ["id", "name"]
	// - f：过滤条件，例如 [["equal", "name", "John Doe"], ["in", "id", "1a", "1b"]]
	// - l：附加子句，例如 "order by id desc limit 20 offset 0"
	// 返回值：
	// - []map[string]interface{}：检索到的记录
	// - error：错误信息
	Get(st string, c []string, f [][]string, l string) ([]map[string]interface{}, error)
	// Update 根据条件修改指定表中的记录。
	// 参数：
	// - st：schema 和表，格式如 "schema.table"
	// - d：要更新的数据
	// - w：WHERE 条件，例如 "id='1a'"
	// 返回值：
	// - error：错误信息
	Update(st string, d map[string]interface{}, w string) error
	// Remove 根据条件删除指定表中的记录。
	// 参数：
	// - st：schema 和表，格式如 "schema.table"
	// - w：WHERE 条件，例如 "id='1a'"
	// 返回值：
	// - error：错误信息
	Remove(st string, w string) error
}

type SharedRepoImpl struct {
	db *sql.DB
}

// NewSharedRepo 创建一个新的 SharedRepoImpl 实例。
// 参数：
// - db：数据库连接
// 返回值：
// - *SharedRepoImpl：SharedRepoImpl 实例
func NewSharedRepo(db *sql.DB) *SharedRepoImpl {
	return &SharedRepoImpl{db: db}
}

// Create 插入一条新记录到指定的表。
// 参数：
// - st：schema 和表，格式如 "schema.table"
// - d：要插入的数据
// 返回值：
// - error：错误信息
func (r *SharedRepoImpl) Create(st string, d map[string]interface{}) error {
	columns, err := get_columns(r.db, st)
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

// Get 根据条件从指定的表中检索记录。
// 参数：
// - st：schema 和表，格式如 "schema.table"
// - c：要检索的列，例如 ["id", "name"]
// - f：过滤条件，例如 [["equal", "name", "John Doe"], ["in", "id", "1a", "1b"]]
// - l：附加子句，例如 "order by id desc limit 20 offset 0"
// 返回值：
// - []map[string]interface{}：检索到的记录
// - error：错误信息
func (r *SharedRepoImpl) Get(st string, c []string, f [][]string, l string) ([]map[string]interface{}, error) {
	if len(c) == 0 {
		var err error
		c, err = get_columns(r.db, st)
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

// Update 根据条件修改指定表中的记录。
// 参数：
// - st：schema 和表，格式如 "schema.table"
// - d：要更新的数据
// - w：WHERE 条件，例如 "id='1a'"
// 返回值：
// - error：错误信息
func (r *SharedRepoImpl) Update(st string, d map[string]interface{}, w string) error {
	columns, err := get_columns(r.db, st)
	if err != nil {
		return err
	}

	q := fmt.Sprintf("update %s set ", st)
	var values []string
	for _, v := range columns {
		if _, ok := d[v]; ok {
			values = append(values, fmt.Sprintf("%s = $%d", v, len(values)+1))
		}
	}
	q += strings.Join(values, ", ")
	q += " where " + w

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
	if err != nil {
		return err
	}

	return nil
}

// Remove 根据条件删除指定表中的记录。
// 参数：
// - st：schema 和表，格式如 "schema.table"
// - w：WHERE 条件，例如 "id='1a'"
// 返回值：
// - error：错误信息
func (r *SharedRepoImpl) Remove(st string, w string) error {
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

type SharedRepoImplMySQL struct {
	db *sql.DB
}

// NewSharedRepoMySQL 创建一个新的 SharedRepoImplMySQL 实例。
// 参数：
// - db：数据库连接
// 返回值：
// - *SharedRepoImplMySQL：SharedRepoImplMySQL 实例
func NewSharedRepoMySQL(db *sql.DB) *SharedRepoImplMySQL {
	return &SharedRepoImplMySQL{db: db}
}

// Create 插入一条新记录到指定的表（MySQL）。
// 参数：
// - st：schema 和表，格式如 "schema.table"
// - d：要插入的数据
// 返回值：
// - error：错误信息
func (r *SharedRepoImplMySQL) Create(st string, d map[string]interface{}) error {
	columns, err := get_columns(r.db, st)
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

// Get 根据条件从指定的表中检索记录。
// 参数：
// - st：schema 和表，格式如 "schema.table"
// - c：要检索的列，例如 ["id", "name"]
// - f：过滤条件，例如 [["equal", "name", "John Doe"], ["in", "id", "1a", "1b"]]
// - l：附加子句，例如 "order by id desc limit 20 offset 0"
// 返回值：
// - []map[string]interface{}：检索到的记录
// - error：错误信息
func (r *SharedRepoImplMySQL) Get(st string, c []string, f [][]string, l string) ([]map[string]interface{}, error) {
	if len(c) == 0 {
		var err error
		c, err = get_columns_mysql(r.db, st)
		if err != nil {
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
		case "jsonb-array-contains":
			whereClauses = append(whereClauses, fmt.Sprintf("JSON_CONTAINS(%s, '\"%s\"')", field, condition[2]))
		case "jsonb-object-contains":
			whereClauses = append(whereClauses, fmt.Sprintf("JSON_CONTAINS(%s, ?, '$')", field))
			params = append(params, fmt.Sprintf(`{"%s": "%s"}`, condition[1], condition[2]))
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

// Update 根据条件修改指定表中的记录（MySQL）。
// 参数：
// - st：schema 和表，格式如 "schema.table"
// - d：要更新的数据
// - w：WHERE 条件，例如 "id='1a'"
// 返回值：
// - error：错误信息
func (r *SharedRepoImplMySQL) Update(st string, d map[string]interface{}, w string) error {
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

// Remove 根据条件删除指定表中的记录（MySQL）。
// 参数：
// - st：schema 和表，格式如 "schema.table"
// - w：WHERE 条件，例如 "id='1a'"
// 返回值：
// - error：错误信息
func (r *SharedRepoImplMySQL) Remove(st string, w string) error {
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
