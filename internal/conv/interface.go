package conv

import "reflect"

func ToInterfaces(value interface{}) []interface{} {
	var reflectValue reflect.Value
	if v, ok := value.(reflect.Value); ok {
		reflectValue = v
		value = v.Interface()
	} else {
		reflectValue = reflect.ValueOf(value)
	}
	reflectKind := reflectValue.Kind()
	for reflectKind == reflect.Ptr {
		reflectValue = reflectValue.Elem()
		reflectKind = reflectValue.Kind()
	}
	length := reflectValue.Len()
	if length == 0 {
		return nil
	}
	array := make([]interface{}, length)
	for i := 0; i < length; i++ {
		array[i] = MapSearch(reflectValue.Index(i).Interface())
	}
	return array
}
