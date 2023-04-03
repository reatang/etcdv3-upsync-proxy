package xnet

import "net"

var privateIPBlocks = []*net.IPNet{
	// 10.0.0.0/8
	{IP: net.ParseIP("10.0.0.0"), Mask: net.CIDRMask(8, 32)},
	// 172.16.0.0/12
	{IP: net.ParseIP("172.16.0.0"), Mask: net.CIDRMask(12, 32)},
	// 192.168.0.0/16
	{IP: net.ParseIP("192.168.0.0"), Mask: net.CIDRMask(16, 32)},
}

func IsPrivateIP(ipAddr string) bool {
	ip := net.ParseIP(ipAddr)

	for _, block := range privateIPBlocks {
		if block.Contains(ip) {
			return true
		}
	}

	return false
}
