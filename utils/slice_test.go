package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInSlice(t *testing.T) {
	testSlice := []string{"a", "ba", "cba"}

	b, _ := InSlice("a", testSlice)
	require.True(t, b)

	b, _ = InSlice("ba", testSlice)
	require.True(t, b)

	b, _ = InSlice("cba", testSlice)
	require.True(t, b)

	b, _ = InSlice("b", testSlice)
	require.False(t, b)
}
