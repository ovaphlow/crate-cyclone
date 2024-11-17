package dbutil

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"
)

// 指定表的列名
// db 数据库连接
// sat schema and table example:"public.setting"
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
	Create(st string, d map[string]interface{}) error
	Get(st string, f [][]string, l string) ([]map[string]interface{}, error)
	Update(st string, d map[string]interface{}, w string) error
	Remove(st string, w string) error
}

type SharedRepoImpl struct {
	db *sql.DB
}

func NewSharedRepo(db *sql.DB) *SharedRepoImpl {
	return &SharedRepoImpl{db: db}
}

// 新增数据
// st schema and table example:"public.setting"
// d 数据
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

// 查询
// st schema and table example:"public.setting"
// f 查询条件 example:[["equal", "name", "John Doe"], ["in", "id", "1a", "1b"]]
// l strings to append example:"order by id desc limit 20 offset 0"
func (r *SharedRepoImpl) Get(st string, f [][]string, l string) ([]map[string]interface{}, error) {
	return nil, nil
}

// 修改
// st schema and table example:"public.setting"
// d 数据
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

	log.Println(q)

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

// 删除
// st schema and table example:"public.setting"
// w strings to append example:"id='1a'"
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
