package yaml

import (
	"encoding/json"
	"gojson/internal/conv"
	"gopkg.in/yaml.v3"
)

func ToJson(content []byte) ([]byte, error) {
	var (
		err    error
		result interface{}
	)
	if result, err = decode(content); err != nil {
		return nil, err
	}
	return json.Marshal(result)
}

func decode(data []byte) (interface{}, error) {
	var (
		result map[string]interface{}
		err    error
	)
	if err = yaml.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return conv.MapSearch(result), nil
}
