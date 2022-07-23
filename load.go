package gojson

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gojson/internal/regex"
)

func (j *Json) parseContent(content []byte, options Options) *Json {
	if options.ContentType == "" {
		options.ContentType = getContentType(content)
		if options.ContentType == "" {
			fmt.Printf("%v,err: %v\n", createErr, invalidContentType)
			j.isValid = false
			return j
		}
	}
	switch options.ContentType {
	case ContentTypeJson:

	case ContentTypeXml:

	case ContentTypeYaml:

	case ContentTypeToml:

	case ContentTypeIni:

	default:
	}

	// 使用json decoder将数据解码成map[string]interface{}形式

	var jsonContent interface{}
	decoder := json.NewDecoder(bytes.NewReader(content))
	// 解码时是否将数字看作字符
	if options.StrNumber {
		decoder.UseNumber()
	}
	if err := decoder.Decode(&jsonContent); err != nil {
		fmt.Printf("%v, err: %v", decodeErr, err)
		j.isValid = false
		return j
	}
	switch jsonContent.(type) {
	// 解码器没有把数据解析成map[string]interface{}的情况
	case string, []byte:
		fmt.Printf("%v", decoder)
		j.isValid = false
		return j
	}
	// 携带解析完后的jsonContent递归下去
	return j.LoadContentWithOptions(jsonContent, options)
}

// getContentType 通过正则表达式判断数据的格式
func getContentType(content []byte) string {
	if json.Valid(content) {
		return ContentTypeJson
	} else if regex.CheckXml(content) {
		return ContentTypeXml
	} else if regex.CheckYaml(content) {
		return ContentTypeYaml
	} else if regex.CheckToml(content) {
		return ContentTypeToml
	} else if regex.CheckIni(content) {
		return ContentTypeIni
	} else {
		return ""
	}
}
