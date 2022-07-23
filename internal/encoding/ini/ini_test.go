package ini

import (
	"fmt"
	"testing"
)

var iniContent = `

;注释
aa=bb
[addr] 
#注释
ip = 127.0.0.1
port=9001
enable=true
command=/bin/echo "gf=GoFrame"

	[DBINFO]
	type=mysql
	user=root
	password=password
[键]
呵呵=值

`

func TestToJson(t *testing.T) {
	jsonStr, err := ToJson([]byte(iniContent))
	if err != nil {
		t.Errorf("%v", err)
	}
	fmt.Println(jsonStr)
}
