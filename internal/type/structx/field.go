package structx

import (
	"errors"
	"reflect"
)

func (f *Field) IsExported() bool {
	return f.Field.PkgPath == ""
}

func (f *Field) Name() string {
	return f.Field.Name
}

func (f *Field) TagStr() string {
	return string(f.Field.Tag)
}

func (f *Field) IsEmbedded() bool {
	return f.Field.Anonymous
}

func Fields(in FieldsInput) ([]Field, error) {
	var (
		ok          bool
		subFieldMap = make(map[string]struct{})
		gotFields   = make([]Field, 0)
		curFieldMap = make(map[string]Field)
	)
	curFields, err := getFieldValues(in.Data)
	if err != nil {
		return nil, err
	}

	for index := 0; index < len(curFields); index++ {
		field := curFields[index]
		curFieldMap[field.Name()] = field
	}

	for index := 0; index < len(curFields); index++ {
		field := curFields[index]
		// 若匿名字段下的子字段和当前字段的字段名若相同
		// 选择写入当前字段
		// 防止出现相同名字的字段发生错误
		if _, ok = subFieldMap[field.Name()]; ok {
			continue
		}
		if field.IsEmbedded() {
			if in.RecursiveOption != RecursiveOptionNone {
				switch in.RecursiveOption {
				case RecursiveOptionEmbeddedNoTag:
					if field.TagStr() != "" {
						break
					}
					fallthrough
				case RecursiveOptionEmbedded:
					subFields, err := Fields(FieldsInput{
						Data:            field.Value,
						RecursiveOption: in.RecursiveOption,
					})
					if err != nil {
						return nil, err
					}
					for i := 0; i < len(subFields); i++ {
						var (
							subField     = subFields[i]
							subFieldName = subField.Name()
						)
						if _, ok = subFieldMap[subFieldName]; ok {
							continue
						}
						subFieldMap[subFieldName] = struct{}{}
						if v, ok := curFieldMap[subFieldName]; !ok {
							//
							gotFields = append(gotFields, subField)
						} else {
							gotFields = append(gotFields, v)
						}
					}
					continue
				}
			}
			continue
		}
		subFieldMap[field.Name()] = struct{}{}
		gotFields = append(gotFields, field)
	}
	return gotFields, nil
}

func getFieldValues(value interface{}) ([]Field, error) {
	var (
		reflectValue reflect.Value
		reflectKind  reflect.Kind
	)
	if v, ok := value.(reflect.Value); ok {
		reflectValue = v
		reflectKind = reflectValue.Kind()
	} else {
		reflectValue = reflect.ValueOf(value)
		reflectKind = reflectValue.Kind()
	}

	// 找到真实数据
	flag := false
	for {
		switch reflectKind {
		case reflect.Ptr:
			if !reflectValue.IsValid() || reflectValue.IsNil() {
				// 如果指针是 *struct 或者 nil,创建一个临时struct
				reflectValue = reflect.New(reflectValue.Type().Elem()).Elem()
				reflectKind = reflectValue.Kind()
			} else {
				reflectValue = reflectValue.Elem()
				reflectKind = reflectValue.Kind()
			}
		case reflect.Array, reflect.Slice:
			reflectValue = reflect.New(reflectValue.Type().Elem()).Elem()
			reflectKind = reflectValue.Kind()
		default:
			flag = true
		}
		if flag {
			break
		}
	}

	if reflectKind != reflect.Struct {
		return nil, errors.New(invalidStructType)
	}
	var (
		fieldType = reflectValue.Type()
		length    = reflectValue.NumField()
		fields    = make([]Field, length)
	)
	for i := 0; i < length; i++ {
		fields[i] = Field{
			Value: reflectValue.Field(i),
			Field: fieldType.Field(i),
		}
	}
	return fields, nil
}
