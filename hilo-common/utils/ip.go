package utils

import (
	"fmt"
	"net"
)

func GetClientIp() (string, error) {
	addrs, err := net.InterfaceAddrs()

	if err != nil {
		return "", err
	}

	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipNet, ok := address.(*net.IPNet); ok && ipNet.IP.IsGlobalUnicast() {
			//if ipNet, ok := address.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String(), nil
			}

		}
	}
	return "", fmt.Errorf("can not find the client ip address")
}
