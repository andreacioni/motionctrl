package utils

import "testing"
import "github.com/stretchr/testify/require"

func TestFromTextSize(t *testing.T) {
	val, err := FromTextSize("121B")

	require.NoError(t, err)
	require.Equal(t, uint64(121), val)

	val, err = FromTextSize("122KB")

	require.NoError(t, err)
	require.Equal(t, uint64(122000), val)

	val, err = FromTextSize("123MB")

	require.NoError(t, err)
	require.Equal(t, uint64(123000000), val)

	val, err = FromTextSize("124GB")

	require.NoError(t, err)
	require.Equal(t, uint64(124000000000), val)

	_, err = FromTextSize("125")

	require.Error(t, err)

	_, err = FromTextSize("0125M")

	require.Error(t, err)

}
