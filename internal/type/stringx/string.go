package stringx

import (
	"fmt"
	"strconv"
	"strings"
)

func IsIndex(in string) bool {
	return strings.HasPrefix(in, "[") && strings.HasSuffix(in, "]")
}

func GetIndex(in string) (int, error) {
	if !IsIndex(in) {
		return -1, fmt.Errorf("invalid index")
	}
	is := strings.TrimLeft(in, "[")
	is = strings.TrimRight(is, "]")
	oint, err := strconv.Atoi(is)
	if err != nil {
		return -1, err
	}
	return oint, nil
}
