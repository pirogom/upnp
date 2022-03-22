package upnp

import (
	"fmt"
	"net"
	"strings"
)

/**
*	GetLocalIntenetIp
**/
func GetLocalIntenetIp() string {
	// PIROGOM
	conn, err := net.Dial("udp", "naver.com:80")
	if err != nil {
		fmt.Println("connection failed")
		return ""
	}
	defer conn.Close()

	return strings.Split(conn.LocalAddr().String(), ":")[0]
}

// This returns the list of local ip addresses which other hosts can connect
// to (NOTE: Loopback ip is ignored).
func GetLocalIPs() ([]*net.IP, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}

	ips := make([]*net.IP, 0)
	for _, addr := range addrs {
		ipnet, ok := addr.(*net.IPNet)
		if !ok {
			continue
		}

		if ipnet.IP.IsLoopback() {
			continue
		}

		ips = append(ips, &ipnet.IP)
	}

	return ips, nil
}
