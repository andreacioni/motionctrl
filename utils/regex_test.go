package utils

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

func TestRegexMustMatch(t *testing.T) {
	if !RegexMustMatch("Camera [0-9]+ Detection status (ACTIVE|PAUSE)", "Camera 0 Detection status ACTIVE") {
		t.Error("not match")
	}

}

func TestRegexFirstSubmatchString(t *testing.T) {
	require.Equal(t, RegexFirstSubmatchString("Camera [0-9]+ Detection status (ACTIVE|PAUSE)", "Camera 0 Detection status ACTIVE"), "ACTIVE")

	require.Equal(t, RegexFirstSubmatchString("Camera [0-9]+ Detection status (ACTIVE|PAUSE)", "Camera 0 Detection status PAUSE"), "PAUSE")

}

func TestRegexConfigList(t *testing.T) {
	testString := "#comment here\n;comment here\nhello 12\nword 11\nnullparam (null)"

	testMap := RegexSubmatchTypedMap("(?m)^([^;#][a-zA-Z0-9_]+) ([a-zA-Z0-9_()]+)$", testString, nil)

	assert.Equal(t, 3, len(testMap))
	assert.Equal(t, "12", testMap["hello"])
	assert.Equal(t, "11", testMap["word"])
	assert.Equal(t, "(null)", testMap["nullparam"])
}
