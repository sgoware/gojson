package gojson

import (
	"gojson/internal/mutex"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		want *Json
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

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
				jsonContent: tt.fields.jsonContent,
				isValid:     tt.fields.isValid,
			}
			if got := j.LoadContent(tt.args.data); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Json.LoadContent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJson_LoadContentWithOptions(t *testing.T) {
	type fields struct {
		mu          *mutex.RWMutex
		jsonContent *interface{}
		isValid     bool
	}
	type args struct {
		data    interface{}
		options Options
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Json
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &Json{
				mu:          tt.fields.mu,
				jsonContent: tt.fields.jsonContent,
				isValid:     tt.fields.isValid,
			}
			if got := j.LoadContentWithOptions(tt.args.data, tt.args.options); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Json.LoadContentWithOptions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJson_LoadHttpResponseBodyWithOptions(t *testing.T) {
	type fields struct {
		mu          *mutex.RWMutex
		jsonContent *interface{}
		isValid     bool
	}
	tests := []struct {
		name   string
		fields fields
		want   *Json
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &Json{
				mu:          tt.fields.mu,
				jsonContent: tt.fields.jsonContent,
				isValid:     tt.fields.isValid,
			}
			if got := j.LoadHttpResponseBodyWithOptions(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Json.LoadHttpResponseBodyWithOptions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJson_Unmarshal(t *testing.T) {
	type fields struct {
		mu          *mutex.RWMutex
		jsonContent *interface{}
		isValid     bool
	}
	type args struct {
		dest interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &Json{
				mu:          tt.fields.mu,
				jsonContent: tt.fields.jsonContent,
				isValid:     tt.fields.isValid,
			}
			if err := j.Unmarshal(tt.args.dest); (err != nil) != tt.wantErr {
				t.Errorf("Json.Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestJson_Get(t *testing.T) {
	type fields struct {
		mu          *mutex.RWMutex
		jsonContent *interface{}
		isValid     bool
	}
	tests := []struct {
		name   string
		fields fields
		want   *Json
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &Json{
				mu:          tt.fields.mu,
				jsonContent: tt.fields.jsonContent,
				isValid:     tt.fields.isValid,
			}
			if got := j.Get(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Json.Get() = %v, want %v", got, tt.want)
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
				jsonContent: tt.fields.jsonContent,
				isValid:     tt.fields.isValid,
			}
			if got := j.LoadFile(tt.args.path); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Json.LoadFile() = %v, want %v", got, tt.want)
			}
		})
	}
}
