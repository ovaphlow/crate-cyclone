package main

import (
	"fmt"
	"strings"
)

type RetrieveOption struct {
	Skip int64
	Take int
}

type RetrieveFilter struct {
	Equal         []string
	ObjectContain []string
	ArrayContain  []string
	Like          []string
	ObjectLike    []string
	In            []string
	Lesser        []string
	Greater       []string
}

func equalBuilder(equal []string) ([]string, []string) {
	if len(equal)%2 != 0 {
		slogger.Debug("equal length is not even")
		return []string{}, []string{}
	}
	var conditions []string
	var params []string
	for i := 0; i < len(equal); i += 2 {
		conditions = append(conditions, fmt.Sprintf("%s = ?", equal[i]))
		params = append(params, equal[i+1])
	}
	return conditions, params
}

func objectContainBuilder(objectContain []string) ([]string, []string) {
	if len(objectContain)%3 != 0 {
		slogger.Debug("objectContain length is not multiple of 3")
		return []string{}, []string{}
	}
	var conditions []string
	var params []string
	for i := 0; i < len(objectContain); i += 3 {
		conditions = append(
			conditions,
			fmt.Sprintf(
				"json_contains(%s, json_object('%s', ?))",
				objectContain[i],
				objectContain[i+1],
			),
		)
		params = append(params, objectContain[i+2])
	}
	return conditions, params
}

func arrayContainBuilder(arrayContain []string) ([]string, []string) {
	if len(arrayContain)%2 != 0 {
		slogger.Debug("arrayContain length is not even")
		return []string{}, []string{}
	}
	var conditions []string
	var params []string
	for i := 0; i < len(arrayContain); i += 2 {
		conditions = append(
			conditions,
			fmt.Sprintf("json_contains(%s, json_array(?))", arrayContain[i]),
		)
		params = append(params, arrayContain[i+1])
	}
	return conditions, params
}

func likeBuilder(like []string) ([]string, []string) {
	if len(like)%2 != 0 {
		slogger.Debug("like length is not even")
		return []string{}, []string{}
	}
	var conditions []string
	var params []string
	for i := 0; i < len(like); i += 2 {
		conditions = append(
			conditions,
			fmt.Sprintf("position(? in %s)", like[i]),
		)
		params = append(params, like[i+1])
	}
	return conditions, params
}

func objectLikeBuilder(objectLike []string) ([]string, []string) {
	if len(objectLike)%3 != 0 {
		slogger.Debug("objectLike length is not multiple of 3")
		return []string{}, []string{}
	}
	var conditions []string
	var params []string
	for i := 0; i < len(objectLike); i += 3 {
		conditions = append(
			conditions,
			fmt.Sprintf(
				"position(? in %s->>'$.%s')",
				objectLike[i],
				objectLike[i+1],
			),
		)
		params = append(params, objectLike[i+2])
	}
	return conditions, params
}

func inBuilder(in []string) ([]string, []string) {
	if len(in) < 2 {
		slogger.Debug("in length is less than 2")
		return []string{}, []string{}
	}
	c := make([]string, len(in)-1)
	for i := range c {
		c[i] = "?"
	}
	return []string{fmt.Sprintf("%s in (%s)", in[0], strings.Join(c, ", "))}, in[1:]
}

func lesserBuilder(lesser []string) ([]string, []string) {
	if len(lesser)%2 != 0 {
		slogger.Debug("lesser length is not even")
		return []string{}, []string{}
	}
	var conditions []string
	var params []string
	for i := 0; i < len(lesser); i += 2 {
		conditions = append(
			conditions,
			fmt.Sprintf("%s <= ?", lesser[i]),
		)
		params = append(params, lesser[i+1])
	}
	return conditions, params
}

func greaterBuilder(greater []string) ([]string, []string) {
	if len(greater)%2 != 0 {
		slogger.Debug("greater length is not even")
		return []string{}, []string{}
	}
	var conditions []string
	var params []string
	for i := 0; i < len(greater); i += 2 {
		conditions = append(
			conditions,
			fmt.Sprintf("%s >= ?", greater[i]),
		)
		params = append(params, greater[i+1])
	}
	return conditions, params
}

func where(conditions []string) string {
	if len(conditions) == 0 {
		return ""
	}
	return fmt.Sprintf(" where %s", strings.Join(conditions, " and "))
}
