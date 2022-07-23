package xml

import (
	"fmt"
	"testing"
)

var (
	xmlContent = `
<?xml version="1.0" encoding="UTF-8"?><doc><username>johngcn</username><password1>123456</password1><password2>123456</password2></doc>
`
)

func TestToJson(t *testing.T) {
	jsonStr, err := ToJson([]byte(xmlContent))
	if err != nil {
		t.Error(err)
	}
	fmt.Println(jsonStr)
}
