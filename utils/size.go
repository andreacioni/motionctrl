package utils

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func FromTextSize(size string) (int64, error) {
	size = strings.Replace(strings.TrimSpace(strings.ToUpper(size)), " ", "", -1)

	re := regexp.MustCompile("^([1-9][0-9]*)(B|KB|MB|GB)$")
	strs := re.FindStringSubmatch(size)

	typeMap := map[string]int64{
		"B":  1,
		"KB": 1000,
		"MB": 1000000,
		"GB": 1000000000,
	}

	if len(strs) != 3 {
		return 0, fmt.Errorf("text doesn't contain a valid size")
	}

	num, err := strconv.ParseInt(strs[1], 10, 64)

	if err != nil {
		return 0, fmt.Errorf("error converting string to int: %s", err)
	}

	mul := typeMap[strs[2]]

	return int64(num) * mul, nil
}
