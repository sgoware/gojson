package conv

import (
	"reflect"
)

// Map
// 不带递归的map(结构体没有嵌套的时候使用)
func Map(value interface{}, tags ...string) map[string]interface{} {
	return convert(value, false, tags...)
}

// MapSearch
// 带递归的map搜索(含嵌套或者嵌套指针结构体的时候使用)
func MapSearch(value interface{}, tags ...string) map[string]interface{} {
	return convert(value, true, tags...)
}

// convert
// 将数据转换成map[string]interface{}
func convert(value interface{}, recursive bool, tags ...string) map[string]interface{} {
	if value == nil {
		return nil
	}
	tagPriority := StructTagPriority
	if len(tags) != 0 {
		tagPriority = append(tags, StructTagPriority...)
	}
	var reflectValue reflect.Value
	if v, ok := value.(reflect.Value); ok {
		reflectValue = v
	} else {
		reflectValue = reflect.ValueOf(value)
	}
	reflectKind := reflectValue.Kind()
	// 是指针的情况:
	//   找到指针所指的真实类型
	for reflectKind == reflect.Ptr {
		reflectValue = reflectValue.Elem()
		reflectKind = reflectValue.Kind()
	}
	switch reflectKind {
	case reflect.Slice, reflect.Array:
		// TODO: 待支持
	case reflect.Map, reflect.Struct, reflect.Interface:
		convertedValue := convertForRecursiveDataStructure(true, value, recursive, tagPriority...)
		if m, ok := convertedValue.(map[string]interface{}); ok {
			return m
		}
		return nil
	default:
		return nil
	}
	return nil
}

func convertForRecursiveDataStructure(isRoot bool, value interface{}, recursive bool, tags ...string) interface{} {
	if isRoot == false && recursive == false {
		return value
	}
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
	switch reflectKind {
	case reflect.Map:
		var (
			mapKeys = reflectValue.MapKeys()
			dataMap = make(map[string]interface{})
		)
		for _, k := range mapKeys {
			dataMap[String(k.Interface())] = convertForRecursiveDataStructure(
				false,
				reflectValue.MapIndex(k).Interface(),
				recursive,
				tags...,
			)
		}
		return dataMap

	case reflect.Struct:
		var dataMap = make(map[string]interface{})
		var (
			rtField     reflect.StructField
			rvField     reflect.Value
			reflectType = reflectValue.Type()
			mapKey      = ""
		)
		numField := reflectValue.NumField()
		for i := 0; i < numField; i++ {
			rtField = reflectType.Field(i)
			rvField = reflectValue.Field(i)
			// 只转换可导出字段
			if !rtField.IsExported() {
				continue
			}
			mapKey = ""
			fieldTag := rtField.Tag

			for _, tag := range tags {
				if mapKey = fieldTag.Get(tag); mapKey != "" {
					break
				}
			}

			if recursive || rtField.Anonymous {
				// 下面开始递归搜索
				var (
					rvAttrField = rvField
					rvAttrKind  = rvField.Kind()
				)
				if rvAttrKind == reflect.Ptr {
					rvAttrField = rvField.Elem()
					rvAttrKind = rvAttrField.Kind()
				}
				switch rvAttrKind {
				case reflect.Struct:
					// 如果结构体没有字段,直接跳过
					if rvAttrField.Type().NumField() == 0 {
						continue
					}
					var (
						rvAttrInterface = rvAttrField.Interface()
					)
					// 如果不是匿名字段
					if !rtField.Anonymous {
						dataMap[mapKey] = convertForRecursiveDataStructure(false, rvAttrInterface, recursive, tags...)
					} else {
						// 如果是匿名字段
						// TODO: 待支持
					}

				case reflect.Array, reflect.Slice:
					// TODO: 待支持
				case reflect.Map:
					// TODO: 待支持
				default:
					// 递归到最底层了
					if rvField.IsValid() {
						dataMap[mapKey] = reflectValue.Field(i).Interface()
					} else {
						dataMap[mapKey] = nil
					}
				}
			} else {
				// 没有开启递归的时候执行的代码
				if rvField.IsValid() {
					dataMap[mapKey] = reflectValue.Field(i).Interface()
				} else {
					dataMap[mapKey] = nil
				}
			}
		}
		if len(dataMap) == 0 {
			return value
		}
		return dataMap

	case reflect.Array, reflect.Slice:
		// TODO: 待支持
	}
	return value
}
