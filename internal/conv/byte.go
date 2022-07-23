package conv

func Bytes(any interface{}) []byte {
	if any == nil {
		return nil
	}
	switch value := any.(type) {
	case string:
		return []byte(value)
	case []byte:
		return value
	default:
		// TODO: 待修改
		return nil
	}
}
