package api

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEmptyAppend(t *testing.T) {
	var a []int

	require.Nil(t, a)

	require.NotEmpty(t, append(a, 1))

	require.ElementsMatch(t, []int{1}, append(a, 1))
}
