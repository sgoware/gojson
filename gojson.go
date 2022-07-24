package gojson

import (
	"encoding/json"
	"errors"
	"fmt"
	"gojson/internal/conv"
	"gojson/internal/mutex"
	"reflect"
)

type Json struct {
	mu          *mutex.RWMutex // 开启安全模式:有指针,关闭时:空指针
	jsonContent *interface{}   // 使用指针传递,效率更高
	isValid     bool           // 查看Json对象是否有效
}

func New() *Json {
	j := &Json{isValid: true} // 默认为有效对象,后续遇到错误设置为无效对象
	return j
}

func (j *Json) LoadContent(data interface{}) *Json {
	nilOption := Options{}
	return j.LoadContentWithOptions(data, nilOption)
}

// LoadContentWithOptions
// 目的将data转换成map[string]interface{}或,map[string][]interface{}的形式
// 使其能够递归调用json的数据
func (j *Json) LoadContentWithOptions(data interface{}, options Options) *Json {
	if data == nil {
		fmt.Printf("%v,err: %v\n", createErr, emptyContest)
		// TODO: 后面用json的时候需要判断Json对象是否有效
		j.isValid = false
		return j
	}
	switch data.(type) {
	// 传入的已经是解码好的json数据的情况
	case map[string]interface{}, map[string][]interface{}:
		j.jsonContent = &data
	// 传入的是字符串或者bytes的情况:
	// 判断数据的格式(json,yaml,toml...),转化成json格式
	// 然后将数据解码成map[string]interface{}的形式
	case string, []byte:
		content := conv.ToBytes(data)
		if len(content) == 0 {
			j.isValid = false
			return j
		}
		return j.parseContent(content, options)
	default:
		var pointedData interface{}
		switch reflect.ValueOf(data).Kind() {

		case reflect.Struct, reflect.Map:
			// 传入的是可递归结构的情况:

			// 如果结构体是接口的情况:
			// 取值然后再递归下去
			// 方法①:
			//   先将结构体marshal成[]bytes,然后unmarshal成map[string]interface{}
			//   但是效率慢,这里想自己写个递归方法
			// 方法②:
			//   直接将结构体转化成map[string]interface{}
			//   利用反射层层递归
			// 这里采用方法②
			pointedData = conv.MapSearch(data, "json")
		case reflect.Slice, reflect.Array:
			// 返回空接口切片
			pointedData = conv.ToInterfaces(data)
		default:
			fmt.Printf("%v, err: %v", createErr, invalidContentType)
			j.isValid = false
			return j
		}
		j.jsonContent = &pointedData
	}
	j.mu = mutex.New(options.Safe) // 创建读写锁
	return j
}

func (j *Json) LoadFileWithOptions() *Json {
	return nil
}

func (j *Json) LoadHttpResponseBodyWithOptions() *Json {
	return nil
}

func (j *Json) Unmarshal(dest interface{}) error {
	if !j.isValid {
		return errors.New(invalidJsonObject)
	}
	j.mu.Lock()
	bytes, err := json.Marshal(*j.jsonContent)
	j.mu.Unlock()
	if err != nil {
		fmt.Println(err)
		return err
	}
	err = json.Unmarshal(bytes, dest)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (j *Json) Get() *Json {
	return nil
}
