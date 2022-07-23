package ini

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"reflect"
	"strings"
)

const (
	sectionNotFound = "section not found"
)

func ToJson(content []byte) ([]byte, error) {
	var (
		err    error
		result interface{}
	)
	if result, err = decode(content); err != nil {
		return nil, err
	}
	return json.Marshal(result)
}

func decode(data []byte) (interface{}, error) {
	var (
		err error
	)
	result := make(map[string]interface{})
	if err = unmarshal(data, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func unmarshal(data []byte, out interface{}) (err error) {
	m := make(map[string]interface{})
	var (
		fieldMap    = make(map[string]interface{})
		bytesReader = bytes.NewReader(data)
		bufioReader = bufio.NewReader(bytesReader)
		section     string
		lastSection string
		haveSection bool
		line        string
	)

	for {
		line, err = bufioReader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		if line = strings.TrimSpace(line); len(line) == 0 {
			continue
		}

		if line[0] == ';' || line[0] == '#' {
			continue
		}
		var (
			sectionBeginPos = strings.Index(line, "[")
			sectionEndPos   = strings.Index(line, "]")
		)
		if sectionBeginPos >= 0 && sectionEndPos >= 2 {
			section = line[sectionBeginPos+1 : sectionEndPos]
			if lastSection == "" {
				lastSection = section
			} else if lastSection != section {
				lastSection = section
				fieldMap = make(map[string]interface{})
			}
			haveSection = true
		} else if haveSection == false {
			continue
		}

		if strings.Contains(line, "=") && haveSection {
			values := strings.Split(line, "=")
			fieldMap[strings.TrimSpace(values[0])] = strings.TrimSpace(strings.Join(values[1:], "="))
			m[section] = fieldMap
		}
	}
	if haveSection == false {
		return errors.New(sectionNotFound)
	}
	v := reflect.ValueOf(out).Elem()
	v.Set(reflect.ValueOf(m))
	return nil
}
