package gojson

import (
	"gojson/internal/type/stringx"
	"strings"
)

func (j *Json) findContent(pattern string) *interface{} {
	if pattern == "" {
		return nil
	}
	if pattern == "." {
		return j.jsonContent
	}
	pointer := j.jsonContent
	nodes := strings.Split(pattern, ".")
	for _, n := range nodes {
		if stringx.IsIndex(n) {
			// 在数组或切片中寻找

			if arr, ok := (*pointer).([]interface{}); ok {
				i, err := stringx.GetIndex(n)
				if err != nil {
					return nil
				}
				arrLen := len(arr)
				if arrLen == 0 ||
					i > arrLen-1 {
					return nil
				}
				return &arr[i]
			}
		} else {
			// 在map中寻找

			if mp, ok := (*pointer).(map[string]interface{}); ok {
				mapValue, ok := mp[n]
				if !ok {
					return nil
				}
				return &mapValue
			}

			if mp, ok := (*pointer).(map[string][]interface{}); ok {
				var result interface{}
				mapValue, ok := mp[n]
				if !ok {
					return nil
				}
				result = &mapValue
				return &result
			}
		}
	}
	return nil
}
