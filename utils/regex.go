package utils

import (
	"regexp"
)

func RegexFirstSubmatchString(regex string, str string) string {
	re := regexp.MustCompile(regex)
	res := re.FindStringSubmatch(str)

	if len(res) > 1 {
		return res[1]
	}

	return ""
}

func RegexSubmatchTypedMap(regex string, str string, typeMapper func(string) interface{}) map[string]interface{} {
	re := regexp.MustCompile(regex)
	strs := re.FindAllStringSubmatch(str, -1)
	retMap := make(map[string]interface{})

	if typeMapper == nil {
		typeMapper = func(s string) interface{} {
			return s
		}
	}

	for i := 0; i < len(strs); i++ {
		retMap[strs[i][1]] = typeMapper(strs[i][2])
	}

	return retMap
}

func RegexMustMatch(regex string, str string) bool {
	re := regexp.MustCompile(regex)
	return re.MatchString(str)
}
