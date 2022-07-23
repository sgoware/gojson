package regex

import (
	"regexp"
)

const (
	xmlRegex      = `^<.+>[\S\s]+<.+>\s*$`
	yamlErrRegex1 = `[\n\r]*[\s\t\w\-\."]+\s*=\s*"""[\s\S]+"""` // 判断嵌套错误情况
	yamlErrRegex2 = `[\n\r]*[\s\t\w\-\."]+\s*=\s*'''[\s\S]+'''` // 同上
	yamlRegex1    = `^[\n\r]*[\w\-\s\t]+\s*:\s*".+"`
	yamlRegex2    = `^[\n\r]*[\w\-\s\t]+\s*:\s*\w+`
	yamlRegex3    = `[\n\r]+[\w\-\s\t]+\s*:\s*".+"`
	yamlRegex4    = `[\n\r]+[\w\-\s\t]+\s*:\s*\w+`
	tomlErrRegex1 = `^[\s\t\n\r]*;.+`
	tomlErrRegex2 = `[\s\t\n\r]+;.+`
	tomlErrRegex3 = `[\n\r]+[\s\t\w\-]+\.[\s\t\w\-]+\s*=\s*.+`
	tomlRegex1    = `[\n\r]*[\s\t\w\-\."]+\s*=\s*".+"`
	tomlRegex2    = `[\n\r]*[\s\t\w\-\."]+\s*=\s*\w+`
	iniRegex      = `\[[\w\.]+\]`
	iniRegex1     = `[\n\r]*[\s\t\w\-\."]+\s*=\s*".+"`
	iniRegex2     = `[\n\r]*[\s\t\w\-\."]+\s*=\s*\w+`
)

func CheckXml(content []byte) bool {
	if isMatch(xmlRegex, content) {
		return true
	}
	return false
}

func CheckYaml(content []byte) bool {
	if !isMatch(yamlErrRegex1, content) && !isMatch(yamlErrRegex2, content) &&
		(isMatch(yamlRegex1, content) ||
			isMatch(yamlRegex2, content) ||
			isMatch(yamlRegex3, content) ||
			isMatch(yamlRegex4, content)) {
		return true
	}
	return false
}

func CheckToml(content []byte) bool {
	if !isMatch(tomlErrRegex1, content) &&
		!isMatch(tomlErrRegex2, content) &&
		!isMatch(tomlErrRegex3, content) &&
		(isMatch(tomlRegex1, content) || isMatch(tomlRegex2, content)) {
		return true
	}
	return false
}

func CheckIni(content []byte) bool {
	if isMatch(iniRegex, content) &&
		(isMatch(iniRegex1, content) || isMatch(iniRegex2, content)) {
		return true
	}
	return false
}

func isMatch(pattern string, src []byte) bool {
	if r, err := getRegexp(pattern); err == nil {
		return r.Match(src)
	}
	return false
}

func getRegexp(pattern string) (regex *regexp.Regexp, err error) {
	regex, err = regexp.Compile(pattern)
	return
}
