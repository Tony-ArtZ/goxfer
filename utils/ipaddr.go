package utils

import "net"

func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()

	if err != nil {
		panic("Unable to get local IP address")
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}

	panic("No valid local IP address found")
}
