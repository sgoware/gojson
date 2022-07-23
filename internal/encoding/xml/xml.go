package xml

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
)

func ToJson(content []byte) ([]byte, error) {
	var (
		err    error
		result interface{}
	)
	// TODO: xml的解码
	if err = xml.NewDecoder(bytes.NewReader(content)).Decode(&result); err != nil {
		return nil, err
	}
	return json.Marshal(result)
}
