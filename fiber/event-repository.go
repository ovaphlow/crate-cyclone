package main

import (
	"fmt"
	"log"
	"strings"
)

var eventColumns = []string{"id", "relation_id", "reference_id", "tags", "detail", "time"}

func EventDefaultFilter(option RetrieveOption, filter RetrieveFilter) ([]Event, error) {
	q := fmt.Sprintf(`select %s from events`, strings.Join(eventColumns, ", "))
	var conditions []string
	var params []string
	if len(filter.Equal) > 0 {
		c, p := equalBuilder(filter.Equal)
		conditions = append(conditions, c...)
		params = append(params, p...)
	}
	if len(filter.ObjectContain) > 0 {
		c, p := objectContainBuilder(filter.ObjectContain)
		conditions = append(conditions, c...)
		params = append(params, p...)
	}
	if len(filter.ArrayContain) > 0 {
		c, p := arrayContainBuilder(filter.ArrayContain)
		conditions = append(conditions, c...)
		params = append(params, p...)
	}
	if len(filter.Like) > 0 {
		c, p := likeBuilder(filter.Like)
		conditions = append(conditions, c...)
		params = append(params, p...)
	}
	if len(filter.ObjectLike) > 0 {
		c, p := objectLikeBuilder(filter.ObjectLike)
		conditions = append(conditions, c...)
		params = append(params, p...)
	}
	if len(filter.In) > 0 {
		c, p := inBuilder(filter.In)
		conditions = append(conditions, c...)
		params = append(params, p...)
	}
	if len(filter.Lesser) > 0 {
		c, p := lesserBuilder(filter.Lesser)
		conditions = append(conditions, c...)
		params = append(params, p...)
	}
	if len(filter.Greater) > 0 {
		c, p := greaterBuilder(filter.Greater)
		conditions = append(conditions, c...)
		params = append(params, p...)
	}
	q += where(conditions)
	q = fmt.Sprintf(`%s order by id desc limit %d, %d`, q, option.Skip, option.Take)
	slogger.Info(q)
	statement, err := MySQL.Prepare(q)
	if err != nil {
		return nil, err
	}
	defer statement.Close()
	var params_ []interface{}
	for _, it := range params {
		params_ = append(params_, it)
	}
	rows, err := statement.Query(params_...)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer rows.Close()
	var result []Event
	for rows.Next() {
		var row Event
		err = rows.Scan(&row.Id, &row.RelationId, &row.ReferenceId, &row.Tags, &row.Detail, &row.Time)
		if err != nil {
			log.Println(err.Error())
			return nil, err
		}
		result = append(result, row)
	}
	if len(result) == 0 {
		return []Event{}, nil
	}
	return result, nil
}
