package utility

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func GetHTTPErrorResponse(title string, status int, r *http.Request) string {
	res := map[string]interface{}{
		"type":     "about:blank",
		"title":    title,
		"status":   status,
		"detail":   "",
		"instance": r.Method + " " + r.RequestURI,
	}
	result, err := json.Marshal(res)
	if err != nil {
		log.Println("Error marshalling HTTP error response:", err)
		return ""
	}
	return string(result)
}

func parseQueryString(filter []string) ([]string, error) {
	if filter[0] == "equal" {
		c, err := strconv.Atoi(filter[1])
		if err != nil {
			return nil, err
		}
		v := filter[2 : 2+c]
		return append([]string{"equal"}, v...), nil
	} else if filter[0] == "in" {
		c, err := strconv.Atoi(filter[1])
		if err != nil {
			return nil, err
		}
		v := filter[2 : 2+c]
		return append([]string{"in"}, v...), nil
	} else if filter[0] == "array-contain" {
		c, err := strconv.Atoi(filter[1])
		if err != nil {
			return nil, err
		}
		v := filter[2 : 2+c]
		return append([]string{"array-contain"}, v...), nil
	} else if filter[0] == "object-contain" {
		c, err := strconv.Atoi(filter[1])
		if err != nil {
			return nil, err
		}
		v := filter[2 : 2+c]
		return append([]string{"object-contain"}, v...), nil
	}
	return nil, nil
}

func ParseQueryString2DefaultFilter(qs string) ([][]string, error) {
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
		parameter, err := parseQueryString(p)
		if err != nil {
			return nil, err
		}
		result = append(result, parameter)
		filter = filter[2+qty:]
	}
	return result, nil
}
