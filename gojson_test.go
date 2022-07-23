package gojson

import (
	"testing"
)

func TestContent(t *testing.T) {
	type Node struct {
		Node      *Node  `json:"node"`
		Name      string `json:"name"`
		Address   string `json:"address,options=[beijing,shanghai]"`
		Age       int    `json:"age"`
		Anonymous bool
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
						Anonymous: false,
					},
					Name:      "json",
					Address:   "beijing",
					Age:       18,
					Anonymous: false,
				},
				Name:      "app",
				Address:   "shanghai",
				Age:       10,
				Anonymous: false,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			New().LoadContent(tt.data)

			//if got := New(); !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("New() = %v, want %v", got, tt.want)
			//}
		})
	}
}
