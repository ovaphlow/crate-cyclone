package utility

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// CreateHTTPResponseRFC9457 创建符合RFC9457格式的HTTP响应。
//
// 参数:
//   - title (string): 响应的标题。
//   - status (int): HTTP状态码。
//   - r (*http.Request): 与响应关联的HTTP请求。
//
// 返回:
//   - map[string]interface{}: HTTP响应的映射。
func CreateHTTPResponseRFC9457(title string, status int, r *http.Request) map[string]interface{} {
	return map[string]interface{}{
		"type":     "about:blank",
		"title":    title,
		"status":   status,
		"detail":   "",
		"instance": r.Method + " " + r.RequestURI,
	}
}

// parseFilterConditions 解析查询字符串为过滤条件。
//
// 参数:
//   - filter ([]string): 包含过滤条件的切片。
//
// 返回:
//   - ([]string, error): 解析后的过滤条件或解析失败时的错误。
func parseFilterConditions(filter []string) ([][]string, error) {
	if filter[0] == "equal" {
		c, err := strconv.Atoi(filter[1])
		if err != nil {
			return nil, err
		}
		if c%2 != 0 {
			return nil, fmt.Errorf("参数数量错误")
		}

		var result [][]string
		for i := 0; i < c; i += 2 {
			result = append(result, []string{"equal", filter[2+i], filter[3+i]})
		}
		return result, nil
	} else if filter[0] == "in" {
		c, err := strconv.Atoi(filter[1])
		if err != nil {
			return nil, err
		}
		v := filter[2 : 2+c]
		return [][]string{append([]string{"in"}, v...)}, nil
	} else if filter[0] == "like" {
		c, err := strconv.Atoi(filter[1])
		if err != nil {
			return nil, err
		}
		if c%2 != 0 {
			return nil, fmt.Errorf("参数数量错误")
		}

		var result [][]string
		for i := 0; i < c; i += 2 {
			result = append(result, []string{"like", filter[2+i], filter[3+i]})
		}
		return result, nil
	} else if filter[0] == "greater-equal" {
		c, err := strconv.Atoi(filter[1])
		if err != nil {
			return nil, err
		}
		if c%2 != 0 {
			return nil, fmt.Errorf("参数数量错误")
		}

		var result [][]string
		for i := 0; i < c; i += 2 {
			result = append(result, []string{"greater-equal", filter[2+i], filter[3+i]})
		}
		return result, nil
	} else if filter[0] == "less-equal" {
		c, err := strconv.Atoi(filter[1])
		if err != nil {
			return nil, err
		}
		if c%2 != 0 {
			return nil, fmt.Errorf("参数数量错误")
		}

		var result [][]string
		for i := 0; i < c; i += 2 {
			result = append(result, []string{"less-equal", filter[2+i], filter[3+i]})
		}
		return result, nil
	} else if filter[0] == "greater" {
		c, err := strconv.Atoi(filter[1])
		if err != nil {
			return nil, err
		}
		if c%2 != 0 {
			return nil, fmt.Errorf("参数数量错误")
		}

		var result [][]string
		for i := 0; i < c; i += 2 {
			result = append(result, []string{"greater", filter[2+i], filter[3+i]})
		}
		return result, nil
	} else if filter[0] == "less" {
		c, err := strconv.Atoi(filter[1])
		if err != nil {
			return nil, err
		}
		if c%2 != 0 {
			return nil, fmt.Errorf("参数数量错误")
		}

		var result [][]string
		for i := 0; i < c; i += 2 {
			result = append(result, []string{"less", filter[2+i], filter[3+i]})
		}
		return result, nil
	} else if filter[0] == "array-contain" {
		c, err := strconv.Atoi(filter[1])
		if err != nil {
			return nil, err
		}
		v := filter[2 : 2+c]
		return [][]string{append([]string{"json-array-contains"}, v...)}, nil
	} else if filter[0] == "object-contain" {
		c, err := strconv.Atoi(filter[1])
		if err != nil {
			return nil, err
		}
		v := filter[2 : 2+c]
		return [][]string{append([]string{"json-object-contains"}, v...)}, nil
	}
	return nil, nil
}

// ConvertQueryStringToDefaultFilter 将查询字符串解析为默认过滤器。
//
// 参数:
//   - qs (string): 原始查询字符串。
//
// 返回:
//   - ([][]string, error): 解析后的过滤��件切片或解析失败时的错误。
func ConvertQueryStringToDefaultFilter(qs string) ([][]string, error) {
	result := [][]string{}
	if qs == "" {
		return result, nil
	}
	filter := strings.Split(qs, ",")
	for len(filter) > 0 {
		qty, err := strconv.Atoi(filter[1])
		if err != nil {
			return nil, err
		}
		p := filter[0 : 2+qty]
		parameter, err := parseFilterConditions(p)
		if err != nil {
			return nil, err
		}
		result = append(result, parameter...)
		filter = filter[2+qty:]
	}
	return result, nil
}
