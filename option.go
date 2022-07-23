package gojson

// Options
// 把选项结构体从Json对象分离出来,Json携带的数据更少,性能更优
type Options struct {
	Safe        bool   // 需要并发安全时开启,使用读写锁
	ContentType string // 设定数据类型,没有设定需要后面程序来判断
	StrNumber   bool   // 是否将数字判断为字符串
}
