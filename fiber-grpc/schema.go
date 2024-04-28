package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"ovaphlow/cratecyclone/models"
	"strings"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type Column struct {
	ColumnName string `json:"column_name"`
	DataType   string `json:"data_type"`
}

type SchemaRepo interface {
	listSchemas() ([]string, error)
	listTables(schema string) ([]string, error)
	listColumns(schema, table string) ([]Column, error)
	save(schema *string, table *string, data map[string]interface{}) error
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

func (r *SchemaRepoImpl) save(schema *string, table *string, data map[string]interface{}) error {
	columns := []string{}
	flags := []string{}
	params := []interface{}{}
	for column := range data {
		columns = append(columns, column)
		flags = append(flags, fmt.Sprintf("$%d", len(flags)+1))
		params = append(params, data[column])
	}
	_, err := r.db.Exec(
		fmt.Sprintf("insert into %s.%s (%s) values (%s)", *schema, *table, strings.Join(columns, ", "), strings.Join(flags, ", ")),
		params...,
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

func (s *SchemaService) ListColumns(schema, table string) ([]Column, error) {
	return s.repo.listColumns(schema, table)
}

func (s *SchemaService) Save(schema *string, table *string, data map[string]interface{}) error {
	columns, err := s.ListColumns(*schema, *table)
	if err != nil {
		return err
	}
	for _, column := range columns {
		if column.ColumnName == "id" && column.DataType == "bigint" {
			node, err := snowflake.NewNode(1)
			if err != nil {
				return err
			}
			data["id"] = node.Generate().Int64()
		}
		if column.ColumnName == "time" {
			data["time"] = time.Now().Format("2006-01-02 15:04:05")
		}
		if column.ColumnName == "state" && column.DataType == "jsonb" {
			randomUUID, err := uuid.NewRandom()
			if err != nil {
				return err
			}
			state := map[string]interface{}{
				"uuid":       randomUUID,
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
		if _, ok := data[column.ColumnName]; !ok {
			validateData = false
			break
		}
	}
	if !validateData {
		return errors.New("数据不完整")
	}
	return s.repo.save(schema, table, data)
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

	app.Post("/crate-api/:schema/:table", func(c *fiber.Ctx) error {
		schema := c.Params("schema")
		table := c.Params("table")
		data := make(map[string]interface{})
		if err := c.BodyParser(&data); err != nil {
			return c.Status(400).JSON(models.ErrorResponse{
				Type:     "about:blank",
				Status:   400,
				Title:    "参数错误",
				Detail:   err.Error(),
				Instance: c.OriginalURL(),
			})
		}
		err := service.Save(&schema, &table, data)
		if err != nil {
			return c.Status(500).JSON(models.ErrorResponse{
				Type:     "about:blank",
				Status:   500,
				Title:    "服务器错误",
				Detail:   err.Error(),
				Instance: c.OriginalURL(),
			})
		}
		return c.Status(201).JSON(data)
	})
}
