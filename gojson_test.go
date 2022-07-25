package gojson

import (
	"fmt"
	"gojson/internal/mutex"
	"reflect"
	"testing"
)

func TestJson_LoadContent(t *testing.T) {
	type fields struct {
		mu          *mutex.RWMutex
		jsonContent *interface{}
		isValid     bool
	}
	type args struct {
		data interface{}
	}

	type NODE struct {
		Node  *NODE   `json:"node"`
		Int   int     `json:"int"`
		Float float64 `json:"float"`
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Json
	}{
		{
			name: "array",
			fields: fields{
				mu:          nil,
				jsonContent: nil,
				isValid:     false,
			},
			args: args{data: []interface{}{NODE{
				Node: &NODE{
					Node:  &NODE{},
					Int:   3,
					Float: 4.33,
				},
				Int:   2,
				Float: 4.544,
			}, NODE{
				Node: &NODE{
					Node:  &NODE{},
					Int:   3,
					Float: 4.33,
				},
				Int:   2,
				Float: 4.544,
			}}},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &Json{
				mu:          tt.fields.mu,
				JsonContent: tt.fields.jsonContent,
				IsValid:     tt.fields.isValid,
			}
			if got := j.LoadContent(tt.args.data); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Json.LoadContent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJson_LoadFile(t *testing.T) {
	type fields struct {
		mu          *mutex.RWMutex
		jsonContent *interface{}
		isValid     bool
	}
	type args struct {
		path string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Json
	}{
		{
			name: "case",
			fields: fields{
				isValid: true,
			},
			args: args{path: "./test.txt"},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &Json{
				mu:          tt.fields.mu,
				JsonContent: tt.fields.jsonContent,
				IsValid:     tt.fields.isValid,
			}
			if got := j.LoadFile(tt.args.path); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Json.LoadFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJson_Get(t *testing.T) {
	node := NestedStruct{
		Nested: &CommonStruct{
			ExampleInt:     34,
			ExampleFloat64: 213.3,
			ExampleString:  "asdwq",
		},
		ExampleInt:     345,
		ExampleFloat64: 2342.234,
		ExampleString:  "asdfawf",
	}
	j := New().LoadContent(node)
	fmt.Println(j.Get("example_int"))
	fmt.Println(j.Get("nested.example_int"))
}

func TestJson_Set(t *testing.T) {
	node := NestedStruct{
		Nested: &CommonStruct{
			ExampleInt:     34,
			ExampleFloat64: 213.3,
			ExampleString:  "asdwq",
		},
		ExampleInt:     345,
		ExampleFloat64: 2342.234,
		ExampleString:  "asdfawf",
	}
	j := New().LoadContent(node)
	fmt.Println(j.Get("example_int"))
	fmt.Println(j.Get("nested.example_int"))
	j.Set("example_int", 222)
	fmt.Println(j.Get("example_int"))
}

func TestJson_Change(t *testing.T) {
	node := NestedStruct{
		Nested: &CommonStruct{
			ExampleInt:     34,
			ExampleFloat64: 213.3,
			ExampleString:  "asdwq",
		},
		ExampleInt:     345,
		ExampleFloat64: 2342.234,
		ExampleString:  "asdfawf",
	}
	j := New().LoadContent(node)
	fmt.Println(j.Get("example_int"))
	j.Set("example_int", 888)
	fmt.Println(j.Get("example_int"))
}

func TestJson_DumpContent(t *testing.T) {
	node := NestedStruct{
		Nested: &CommonStruct{
			ExampleInt:     34,
			ExampleFloat64: 213.3,
			ExampleString:  "asdwq",
		},
		ExampleInt:     345,
		ExampleFloat64: 2342.234,
		ExampleString:  "asdfawf",
	}
	j := New().LoadContent([]interface{}{node, node, node})
	j.Dump()
}
