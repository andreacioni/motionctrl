package utils

import "testing"

func TestRegexMustMatch(t *testing.T) {
	if !RegexMustMatch("Camera [0-9]+ Detection status (ACTIVE|PAUSE)", "Camera 0 Detection status ACTIVE") {
		t.Error("not match")
	}

}

func TestRegexFirstSubmatchString(t *testing.T) {
	if RegexFirstSubmatchString("Camera [0-9]+ Detection status (ACTIVE|PAUSE)", "Camera 0 Detection status ACTIVE") != "ACTIVE" {
		t.Error("not match")
	}

	if RegexFirstSubmatchString("Camera [0-9]+ Detection status (ACTIVE|PAUSE)", "Camera 0 Detection status PAUSE") != "PAUSE" {
		t.Error("not match")
	}
}

func TestRegexConfigList(t *testing.T) {
	testString := "#comment here\n;comment here\nhello 12\nword 11\nnullparam (null)"

	testMap := RegexSubmatchMap("(?m)^([^;#][a-zA-Z0-9_]+) ([a-zA-Z0-9_()]+)$", testString)

	t.Log(testMap)
	if len(testMap) != 2 || testMap["hello"] != 12 || testMap["word"] != 11 {
		t.Error("not match")
	}
}
