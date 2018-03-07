package utils

import (
	"net"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsLocalIPTrue(t *testing.T) {
	ip := net.ParseIP("127.0.0.1")

	require.NotNil(t, ip)

	res, err := IsLocalIP(ip)

	require.NoError(t, err)
	require.True(t, res)

	ip = net.ParseIP("118.123.123.231")

	require.NotNil(t, ip)

	res, err = IsLocalIP(ip)

	require.NoError(t, err)
	require.False(t, res)

	//Nil
	res, err = IsLocalIP(nil)

	require.NoError(t, err)
	require.False(t, res)
}
