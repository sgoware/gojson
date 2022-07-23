package toml

import (
	"encoding/json"
	"github.com/BurntSushi/toml"
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
	var result interface{}
	if err := toml.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return result, nil
}
