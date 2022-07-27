package gojson

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"gojson/internal/mutex"
	"io"
	"os"
)

type Json struct {
	mu          *mutex.RWMutex // 开启安全模式:有指针,关闭时:空指针
	JsonContent *interface{}   // 使用指针传递,效率更高
}

func New() *Json {
	j := &Json{} // 默认为有效对象,后续遇到错误设置为无效对象
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
		return nil
	}
	content, err := j.convertContent(data, options)
	if err != nil {
		fmt.Printf("%v, err: %v", createErr, err)
		return nil
	}
	j.JsonContent = &content
	j.mu = mutex.New(options.Safe) // 创建读写锁
	return j
}

func (j *Json) LoadFile(path string) *Json {
	nilOption := Options{}
	return j.LoadFileWithOptions(path, nilOption)
}
func (j *Json) LoadFileWithOptions(path string, options Options) *Json {
	var content []byte
	file, err := os.Open(path)
	if err != nil {
		fmt.Printf("%v, err: %v, %v", createErr, readFileErr, err)
		return nil
	}
	r := bufio.NewReader(file)
	for {
		lineBytes, err := r.ReadBytes('\n')
		if err != nil && err != io.EOF {
			fmt.Printf("%v, err: %v, %v", createErr, readFileErr, err)
			return nil
		}
		content = append(content, lineBytes...)
		if err == io.EOF {
			break
		}
	}
	return j.LoadContentWithOptions(content, options)
}

func (j *Json) Unmarshal(dest interface{}) error {
	if j == nil {
		return errors.New(invalidJsonObject)
	}
	j.mu.Lock()
	defer j.mu.Unlock()
	bytes, err := json.Marshal(*j.JsonContent)
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

// Get 输出json字符串指定路径的内容
func (j *Json) Get(pattern string) interface{} {
	if j == nil {
		fmt.Printf("%v, err: %v", getErr, invalidJsonObject)
		return ""
	}
	j.mu.Lock()
	defer j.mu.Unlock()
	pointer := j.findContentPointer(pattern)
	if pointer != nil {
		return *pointer
	}
	fmt.Printf("%v, err: %v", getErr, invalidPattern)
	return nil
}

// Set 支持数据替换,插入,删除  data为空为删除
func (j *Json) Set(pattern string, data interface{}) *Json {
	nilOptions := Options{}
	return j.SetWithOptions(pattern, data, nilOptions)
}

func (j *Json) SetWithOptions(pattern string, data interface{}, options Options) *Json {
	if j == nil {
		fmt.Printf("%v, err: %v", setErr, invalidPattern)
		return nil
	}
	j.mu.Lock()
	defer j.mu.Unlock()
	err := j.setContentWithOptions(pattern, data, options)
	if err != nil {
		fmt.Printf("%v, err: %v", setErr, invalidPattern)
		return nil
	}
	return j
}

func (j *Json) Dump() *Json {
	nilOptions := DumpOption{}
	j.DumpWithOptions(j, nilOptions)
	return j
}

func (j *Json) DumpContent() *Json {
	nilOptions := DumpOption{}
	j.DumpWithOptions(j.JsonContent, nilOptions)
	return j
}

func (j *Json) DumpWithOptions(data interface{}, options DumpOption) *Json {
	j.mu.Lock()
	defer j.mu.Lock()
	if j == nil {
		fmt.Printf("%v, err: %v", dumpErr, invalidContentType)
		return j
	}
	dumpWithOption(data, options)
	return j
}
