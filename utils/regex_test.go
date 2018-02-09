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
