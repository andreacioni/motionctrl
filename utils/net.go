package utils

import "net"

func IsLocalIP(localIP net.IP) (bool, error) {
	ifaces, err := net.Interfaces()

	if err != nil {
		return false, err
	}

	for _, i := range ifaces {
		addrs, err := i.Addrs()

		if err != nil {
			return false, err
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			if ip != nil && ip.Equal(localIP) {
				return true, nil
			}
		}
	}

	return false, nil
}
