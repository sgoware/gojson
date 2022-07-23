package gojson

import (
	"fmt"
	"testing"
)

func TestContent(t *testing.T) {
	type Node struct {
		Node      *Node       `json:"node"`
		Name      string      `json:"name"`
		Address   string      `json:"address"`
		Age       int         `json:"age"`
		Iface     interface{} `json:"iface"`
		anonymous bool
	}
	tests := []struct {
		name string
		data interface{}
		want *Json
	}{
		{
			name: "json",
			data: `{"name":"john", "score":"100"}`,
			want: nil,
		},
		{
			name: "struct",
			data: Node{
				Node: &Node{
					Node: &Node{
						Node:      nil,
						Name:      "asd",
						Address:   "qwe",
						Age:       1,
						Iface:     35,
						anonymous: false,
					},
					Name:      "json",
					Address:   "beijing",
					Age:       18,
					Iface:     234,
					anonymous: false,
				},
				Name:      "app",
				Address:   "shanghai",
				Age:       10,
				Iface:     34,
				anonymous: false,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := New().LoadContent(tt.data)
			fmt.Println(*j.jsonContent)
			jj := Node{}
			j.Unmarshal(&jj)
			fmt.Println(jj)
			//if got := New(); !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("New() = %v, want %v", got, tt.want)
			//}
		})
	}
}
