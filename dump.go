package gojson

import (
	"bytes"
	"fmt"
	"gojson/internal/type/reflection"
	"gojson/internal/type/structx"
	"reflect"
	"strings"
)

type DumpOption struct {
	WithType     bool // 是否输出类型名
	ExportedOnly bool // 是否输出不可导出字段
}

// 递归时保存的数据缓存
type dumpBuf struct {
	Value           interface{}
	ReflectValue    reflect.Value
	ReflectTypeName string // 用于保存类型名
	Option          DumpOption
	Indent          string // 缩进
	NewIndent       string // 新的缩进,每一层增加一个缩进
	Buffer          *bytes.Buffer
}

const (
	jsonIndent = `    ` // json格式的缩进
	yamlIndent = `  `   // yaml格式的缩进
)

var (
	// DefaultTrimChars
	// 对照ASCII码表定义一些空白字符
	// 用于分割字符串
	DefaultTrimChars = string([]byte{
		'\t', // 水平制表符
		'\v', // 垂直制表符
		'\n', // 换行符
		'\r', // 回车符
		'\f', // 换页符
		' ',  // 空格
		0xA0, // &nbsp
		0x00, // 空字符
		0x85, // 删除符
		0x1C, // 文件分隔符
	})
)

func dump(values ...interface{}) {
	for _, value := range values {
		dumpWithOption(value, DumpOption{
			WithType:     false,
			ExportedOnly: false,
		})
	}
}

func dumpWithOption(value interface{}, options DumpOption) {
	buffer := bytes.NewBuffer(nil)
	doDump(value, "", buffer, DumpOption{
		WithType:     options.WithType,
		ExportedOnly: options.ExportedOnly,
	})
	fmt.Println(buffer.String())
}

func doDump(value interface{}, indent string, buffer *bytes.Buffer, option DumpOption) {
	if value == nil {
		buffer.WriteString(`<nil>`)
		return
	}
	var reflectValue reflect.Value
	if v, ok := value.(reflect.Value); ok {
		reflectValue = v
		if v.IsValid() && v.CanInterface() {
			value = v.Interface()
		} else {
			if convertedValue, ok := reflection.ValueToInterface(v); ok {
				value = convertedValue
			}
		}
	} else {
		reflectValue = reflect.ValueOf(value)
	}
	if value == nil {
		buffer.WriteString(`<nil>`)
		return
	}
	var (
		reflectKind     = reflectValue.Kind()
		reflectTypeName = reflectValue.Type().String()
		newIndent       = indent + jsonIndent
	)
	// 统一[]byte类型名
	reflectTypeName = strings.ReplaceAll(reflectTypeName, `[]uint8`, `[]byte`)
	if !option.WithType {
		reflectTypeName = ""
	}
	for reflectKind == reflect.Ptr {
		reflectValue = reflectValue.Elem()
		reflectKind = reflectValue.Kind()
	}
	var (
		newDumpBuf = dumpBuf{
			Value:           value,
			ReflectValue:    reflectValue,
			ReflectTypeName: reflectTypeName,
			Option:          option,
			Indent:          indent,
			NewIndent:       newIndent,
			Buffer:          buffer,
		}
	)
	switch reflectKind {
	case reflect.Slice, reflect.Array:
		doDumpSlice(newDumpBuf)

	case reflect.Map:
		doDumpMap(newDumpBuf)

	case reflect.Struct:
		doDumpStruct(newDumpBuf)

	case reflect.Interface:
		doDump(newDumpBuf.ReflectValue.Elem(), indent, buffer, option)

	case reflect.String:
		doDumpString(newDumpBuf)

	case reflect.Bool:
		doDumpBool(newDumpBuf)

	case
		reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64,
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Float32,
		reflect.Float64,
		reflect.Complex64,
		reflect.Complex128:
		doDumpNumber(newDumpBuf)

	default:
		doDumpDefault(newDumpBuf)
	}
}

func doDumpSlice(buf dumpBuf) {
	if b, ok := buf.Value.([]byte); ok {
		if !buf.Option.WithType {
			buf.Buffer.WriteString(fmt.Sprintf(`"%s"`, replaceSlashesForString(string(b))))
		} else {
			buf.Buffer.WriteString(fmt.Sprintf(
				`%s(%d) "%s"`,
				buf.ReflectTypeName,
				len(string(b)),
				string(b),
			))
		}
		return
	}
	length := buf.ReflectValue.Len()
	if length == 0 {
		// 切片长度为0的情况
		if !buf.Option.WithType {
			buf.Buffer.WriteString("[]")
		} else {
			buf.Buffer.WriteString(fmt.Sprintf("%s(0) []", buf.ReflectTypeName))
		}
		return
	}

	// 打印开头
	if !buf.Option.WithType {
		buf.Buffer.WriteString("[\n")
	} else {
		buf.Buffer.WriteString(fmt.Sprintf("%s(%d) [\n", buf.ReflectTypeName, length))
	}

	for i := 0; i < length; i++ {
		// 遍历切片
		buf.Buffer.WriteString(buf.NewIndent)
		doDump(buf.ReflectValue.Index(i), buf.NewIndent, buf.Buffer, buf.Option)
		buf.Buffer.WriteString(",\n")
	}
	buf.Buffer.WriteString(fmt.Sprintf("%s]", buf.Indent))
}

func doDumpMap(buf dumpBuf) {
	var mapKeys = make([]reflect.Value, 0)
	for _, key := range buf.ReflectValue.MapKeys() {
		// 获取map的所有key
		if !key.CanInterface() {
			// 该字段为不可导出字段直接跳到下一个key
			continue
		}
		mapKey := key
		mapKeys = append(mapKeys, mapKey)
	}
	if len(mapKeys) == 0 {
		// map为空的情况
		if !buf.Option.WithType {
			buf.Buffer.WriteString("{}")
		} else {
			buf.Buffer.WriteString(fmt.Sprintf("%s(0) {}", buf.ReflectTypeName))
		}
		return
	}
	var (
		maxSpaceNum = 0 // 根据key名词的最大长度确定缩进
		tmpSpaceNum = 0
		mapKeyStr   = ""
	)
	for _, key := range mapKeys {
		tmpSpaceNum = len(fmt.Sprintf(`%v`, key.Interface()))
		if tmpSpaceNum > maxSpaceNum {
			maxSpaceNum = tmpSpaceNum
		}
	}
	if !buf.Option.WithType {
		buf.Buffer.WriteString("{\n")
	} else {
		buf.Buffer.WriteString(fmt.Sprintf("%s(%d) {\n", buf.ReflectTypeName, len(mapKeys)))
	}
	for _, mapKey := range mapKeys {
		// 根据key遍历map

		// 输出map的key和缩进
		curSpaceNum := len(fmt.Sprintf(`%v`, mapKey.Interface()))
		if mapKey.Kind() == reflect.String {
			mapKeyStr = fmt.Sprintf(`"%v"`, mapKey.Interface())
		} else {
			mapKeyStr = fmt.Sprintf(`%v`, mapKey.Interface())
		}
		if !buf.Option.WithType {
			buf.Buffer.WriteString(fmt.Sprintf(
				"%s%v:%s",
				buf.NewIndent,
				mapKeyStr,
				strings.Repeat(" ", maxSpaceNum-curSpaceNum+1),
			))
		} else {
			buf.Buffer.WriteString(fmt.Sprintf(
				"%s%s(%v):%s",
				buf.NewIndent,
				mapKey.Type().String(),
				mapKeyStr,
				strings.Repeat(" ", maxSpaceNum-curSpaceNum+1),
			))
		}
		// 输出map的value
		doDump(buf.ReflectValue.MapIndex(mapKey), buf.NewIndent, buf.Buffer, buf.Option)
		buf.Buffer.WriteString(",\n")
	}
	buf.Buffer.WriteString(fmt.Sprintf("%s}", buf.Indent))
}

func doDumpStruct(buf dumpBuf) {
	// 获取结构体所有字段
	structFields, _ := structx.Fields(structx.FieldsInput{
		Data:            buf.Value,
		RecursiveOption: structx.RecursiveOptionEmbedded,
	})

	var (
		maxSpaceNum = 0
		tmpSpaceNum = 0
	)
	for _, field := range structFields {
		if buf.Option.ExportedOnly && !field.IsExported() {
			continue
		}
		tmpSpaceNum = len(field.Name())
		if tmpSpaceNum > maxSpaceNum {
			maxSpaceNum = tmpSpaceNum
		}
	}

	// 打印开头
	if !buf.Option.WithType {
		buf.Buffer.WriteString("{\n")
	} else {
		buf.Buffer.WriteString(fmt.Sprintf("%s(%d) {\n",
			buf.ReflectTypeName,
			len(structFields)),
		)
	}
	for _, field := range structFields {
		// 遍历字段

		if buf.Option.ExportedOnly && !field.IsExported() {
			continue
		}
		curSpaceNum := len(fmt.Sprintf(`%v`, field.Name()))
		buf.Buffer.WriteString(fmt.Sprintf(
			"%s%s:%s",
			buf.NewIndent,
			field.Name(),
			strings.Repeat(" ", maxSpaceNum-curSpaceNum+1),
		))
		// 递归遍历下一个层级
		doDump(field.Value, buf.NewIndent, buf.Buffer, buf.Option)
		buf.Buffer.WriteString(",\n")
	}
	buf.Buffer.WriteString(fmt.Sprintf("%s}", buf.Indent))
}

func doDumpNumber(buf dumpBuf) {
	doDumpDefault(buf)
}

func doDumpString(buf dumpBuf) {
	s := buf.ReflectValue.String()
	if !buf.Option.WithType {
		buf.Buffer.WriteString(fmt.Sprintf(`"%v"`, replaceSlashesForString(s)))
	} else {
		buf.Buffer.WriteString(fmt.Sprintf(
			`%s(%d) "%v"`,
			buf.ReflectTypeName,
			len(s),
			replaceSlashesForString(s),
		))
	}
}

func doDumpBool(buf dumpBuf) {
	var s string
	if buf.ReflectValue.Bool() {
		s = `true`
	} else {
		s = `false`
	}
	if buf.Option.WithType {
		s = fmt.Sprintf(`bool(%s)`, s)
	}
	buf.Buffer.WriteString(s)
}

func doDumpDefault(buf dumpBuf) {
	var s string
	if buf.ReflectValue.IsValid() && buf.ReflectValue.CanInterface() {
		s = fmt.Sprintf("%v", buf.ReflectValue.Interface())
	}
	if s == "" {
		s = fmt.Sprintf("%v", buf.Value)
	}
	s = trim(s, `<>`)
	if !buf.Option.WithType {
		buf.Buffer.WriteString(s)
	} else {
		buf.Buffer.WriteString(fmt.Sprintf("%s(%s)", buf.ReflectTypeName, s))
	}
}

// 将字符串形式的斜杠替换成byte形式的
func replaceSlashesForString(s string) string {
	return replaceByMap(s, map[string]string{
		`"`:  `\"`,
		"\r": `\r`,
		"\t": `\t`,
		"\n": `\n`,
	})
}

func replaceByMap(origin string, replaces map[string]string) string {
	for k, v := range replaces {
		origin = strings.Replace(origin, k, v, -1)
	}
	return origin
}

func trim(str string, characterMask ...string) string {
	trimChars := DefaultTrimChars
	if len(characterMask) > 0 {
		trimChars += characterMask[0]
	}
	return strings.Trim(str, trimChars)
}
