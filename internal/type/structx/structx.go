package structx

import (
	"reflect"
)

type Field struct {
	Value reflect.Value
	Field reflect.StructField
}

// FieldsInput is the input parameter struct type for function Fields.
type FieldsInput struct {
	Data            interface{}
	RecursiveOption int // 对结构体中匿名字段操作的具体选项
}

const (
	invalidStructType            = "given value should be the among type of struct/*struct/[]struct/[]*struct"
	RecursiveOptionNone          = 0 // 对匿名字段不递归
	RecursiveOptionEmbedded      = 1 // 对匿名字段递归
	RecursiveOptionEmbeddedNoTag = 2 // 对匿名字段且字段没有标签时递归
)
