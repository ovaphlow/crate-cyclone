package infrastructure

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

type Option struct {
	Count   bool
	Offset  int64
	Take    int
	OrderBy string
}

type SchemaRepo interface {
	listSchemas() ([]string, error)
	listTables(schema string) ([]string, error)
	listColumns(schema, table string) ([]string, error)
	save(schema, table string, data map[string]interface{}) error
	get(schema, table, id string) (map[string]interface{}, error)
	retrieve(schema, table string, conditions [][]string, option Option) ([]map[string]interface{}, error)
	update(schema, table, id string, data map[string]interface{}) error
	remove(schema, table, id string) error
}

type SchemaRepoImpl struct {
	db *sql.DB
}

func NewSchemaRepoImpl(db *sql.DB) *SchemaRepoImpl {
	return &SchemaRepoImpl{db: db}
}

func (r *SchemaRepoImpl) listSchemas() ([]string, error) {
	schemas := []string{}
	rows, err := r.db.Query("SELECT schema_name FROM information_schema.schemata")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var schema string
		err := rows.Scan(&schema)
		if err != nil {
			return nil, err
		}
		schemas = append(schemas, schema)
	}
	return schemas, nil
}

func (r *SchemaRepoImpl) listTables(schema string) ([]string, error) {
	tables := []string{}
	rows, err := r.db.Query("SELECT table_name FROM information_schema.tables WHERE table_schema = $1", schema)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var table string
		err := rows.Scan(&table)
		if err != nil {
			return nil, err
		}
		tables = append(tables, table)
	}
	return tables, nil
}

func (r *SchemaRepoImpl) listColumns(schema, table string) ([]string, error) {
	columns := []string{}
	rows, err := r.db.Query(
		`
		SELECT column_name FROM information_schema.columns
		WHERE table_schema = $1 AND table_name = $2
		order by ordinal_position asc
		`,
		schema,
		table,
	)
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

func (r *SchemaRepoImpl) save(schema string, table string, data map[string]interface{}) error {
	columns := []string{}
	flags := []string{}
	params := []interface{}{}
	for column := range data {
		columns = append(columns, column)
		flags = append(flags, fmt.Sprintf("$%d", len(flags)+1))
		params = append(params, data[column])
	}
	_, err := r.db.Exec(
		fmt.Sprintf("insert into %s.%s (%s) values (%s)", schema, table, strings.Join(columns, ", "), strings.Join(flags, ", ")),
		params...,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *SchemaRepoImpl) get(schema, table, id string) (map[string]interface{}, error) {
	data := make(map[string]interface{})
	rows, err := r.db.Query(
		fmt.Sprintf("select * from %s.%s where id = $1", schema, table),
		id,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if rows.Next() {
		columns, err := rows.Columns()
		if err != nil {
			return nil, err
		}
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range columns {
			valuePtrs[i] = &values[i]
		}
		err = rows.Scan(valuePtrs...)
		if err != nil {
			return nil, err
		}
		for i, col := range columns {
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				data[col] = string(b)
			} else {
				data[col] = val
			}
		}
	}
	return data, nil
}

func (r *SchemaRepoImpl) retrieve(schema, table string, conditions [][]string, option Option) ([]map[string]interface{}, error) {
	return nil, nil
}

func (r *SchemaRepoImpl) update(schema, table, id string, data map[string]interface{}) error {
	columns := []string{}
	params := []interface{}{}
	for column := range data {
		if data[column] == nil {
			columns = append(columns, column)
		} else {
			columns = append(columns, fmt.Sprintf("%s = $%d", column, len(columns)+1))
			params = append(params, data[column])
		}
	}
	_, err := r.db.Exec(
		fmt.Sprintf("update %s.%s set %s where id = $%d", schema, table, strings.Join(columns, ", "), len(params)+1),
		append(params, id)...,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *SchemaRepoImpl) remove(schema, table, id string) error {
	_, err := r.db.Exec(
		fmt.Sprintf("update %s.%s set state = state || jsonb_build_object('deleted_at', '%s') where id = $1", schema, table, time.Now().Format("2006-01-02 15:04:05")),
		id,
	)
	if err != nil {
		return err
	}
	return nil
}

type SchemaService struct {
	repo SchemaRepo
}

func NewSchemaService(repo SchemaRepo) *SchemaService {
	return &SchemaService{repo: repo}
}

func (s *SchemaService) ListSchemas() ([]string, error) {
	return s.repo.listSchemas()
}

func (s *SchemaService) ListTables(schema string) ([]string, error) {
	return s.repo.listTables(schema)
}

func (s *SchemaService) ListColumns(schema, table string) ([]string, error) {
	return s.repo.listColumns(schema, table)
}

func (s *SchemaService) Save(schema string, table string, data map[string]interface{}) error {
	columns, err := s.ListColumns(schema, table)
	if err != nil {
		return err
	}
	for _, column := range columns {
		if column == "id" {
			data["id"], err = GenerateKsuid()
			if err != nil {
				Slogger.Error(err.Error())
				return err
			}
		}
		if column == "time" {
			data["time"] = time.Now().Format("2006-01-02 15:04:05")
		}
		if column == "state" {
			state := map[string]interface{}{
				"created_at": time.Now().Format("2006-01-02 15:04:05"),
			}
			stateJson, err := json.Marshal(state)
			if err != nil {
				return err
			}
			data["state"] = string(stateJson)
		}
	}
	validateData := true
	for _, column := range columns {
		if _, ok := data[column]; !ok {
			validateData = false
			break
		}
	}
	if !validateData {
		return errors.New("数据不匹配")
	}
	return s.repo.save(schema, table, data)
}

func (s *SchemaService) Get(schema, table, id string) (map[string]interface{}, error) {
	return s.repo.get(schema, table, id)
}

func (s *SchemaService) List(schema, table string, conditions [][]string, option Option) ([]map[string]interface{}, error) {
	return s.repo.retrieve(schema, table, conditions, option)
}

func (s *SchemaService) Update(schema, table, id string, data map[string]interface{}) error {
	data[fmt.Sprintf(
		"state = state || jsonb_build_object('updated_at', '%s')",
		time.Now().Format("2006-01-02 15:04:05"),
	)] = nil
	return s.repo.update(schema, table, id, data)
}

func (s *SchemaService) Remove(schema, table, id string) error {
	return s.repo.remove(schema, table, id)
}
