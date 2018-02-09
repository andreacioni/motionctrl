package utils

import (
	"regexp"
)

func RegexFirstSubmatchString(regex string, str string) string {
	re := regexp.MustCompile(regex)
	return re.FindStringSubmatch(str)[1]
}

func RegexMustMatch(regex string, str string) bool {
	re := regexp.MustCompile(regex)
	return re.MatchString(str)
}
