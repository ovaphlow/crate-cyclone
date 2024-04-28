package main

import (
	"database/sql"
	"ovaphlow/cratecyclone/models"

	"github.com/gofiber/fiber/v2"
)

type Column struct {
	ColumnName string `json:"column_name"`
	DataType   string `json:"data_type"`
}

type SchemaRepo interface {
	listSchemas() ([]string, error)
	listTables(schema string) ([]string, error)
	listColumns(schema, table string) ([]Column, error)
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

func (r *SchemaRepoImpl) listColumns(schema, table string) ([]Column, error) {
	columns := []Column{}
	rows, err := r.db.Query(
		`
		SELECT column_name, data_type FROM information_schema.columns
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
		var column Column
		err := rows.Scan(&column.ColumnName, &column.DataType)
		if err != nil {
			return nil, err
		}
		columns = append(columns, column)
	}
	return columns, nil
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

func (s *SchemaService) ListColumns(schema, table string) ([]Column, error) {
	return s.repo.listColumns(schema, table)
}

func AddSchemaEndpoints(app *fiber.App, service *SchemaService) {
	app.Get("/crate-api/db-schema", func(c *fiber.Ctx) error {
		schemas, err := service.ListSchemas()
		if err != nil {
			return c.Status(500).JSON(models.ErrorResponse{
				Type:     "about:blank",
				Status:   500,
				Title:    "服务器错误",
				Detail:   err.Error(),
				Instance: c.OriginalURL(),
			})
		}
		return c.JSON(schemas)
	})

	app.Get("/crate-api/:schema/db-table", func(c *fiber.Ctx) error {
		schema := c.Params("schema")
		tables, err := service.ListTables(schema)
		if err != nil {
			return c.Status(500).JSON(models.ErrorResponse{
				Type:     "about:blank",
				Status:   500,
				Title:    "服务器错误",
				Detail:   err.Error(),
				Instance: c.OriginalURL(),
			})
		}
		return c.JSON(tables)
	})

	app.Get("/crate-api/:schema/:table", func(c *fiber.Ctx) error {
		schema := c.Params("schema")
		table := c.Params("table")
		columns, err := service.ListColumns(schema, table)
		if err != nil {
			return c.Status(500).JSON(models.ErrorResponse{
				Type:     "about:blank",
				Status:   500,
				Title:    "服务器错误",
				Detail:   err.Error(),
				Instance: c.OriginalURL(),
			})
		}
		return c.JSON(columns)
	})
}
