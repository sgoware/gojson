package conv

import (
	"encoding/json"
	"fmt"
	"reflect"
)

func Bytes(any interface{}) []byte {
	if any == nil {
		return nil
	}
	switch value := any.(type) {
	case string:
		return []byte(value)
	case []byte:
		return value
	default:
		originKind := reflect.ValueOf(any).Kind()
		switch originKind {
		case reflect.Map:
			bytes, err := json.Marshal(any)
			if err != nil {
				fmt.Println("err")
				return nil
			}
			return bytes
		case reflect.Array, reflect.Slice:
			// TODO: 待支持
		}
	}
	return nil
}
