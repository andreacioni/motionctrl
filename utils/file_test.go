package utils

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestListFiles(t *testing.T) {
	_, list, _, err := ListFiles("../", nil)
	require.NoError(t, err)
	require.True(t, len(list) > 0)

	_, list, _, err = ListFiles("../", func(f os.FileInfo) bool { return false })
	require.NoError(t, err)
	require.Equal(t, 0, len(list))
}
