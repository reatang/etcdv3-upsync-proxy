package xnet

import "net"

var privateIPBlocks = []*net.IPNet{
	// A类：10.0.0.0/8
	{IP: net.ParseIP("10.0.0.0"), Mask: net.CIDRMask(8, 32)},
	// B类：172.16.0.0/12
	{IP: net.ParseIP("172.16.0.0"), Mask: net.CIDRMask(12, 32)},
	// C类：192.168.0.0/16
	{IP: net.ParseIP("192.168.0.0"), Mask: net.CIDRMask(16, 32)},
}

// IsPrivateIP 计算私有IP
func IsPrivateIP(ipAddr string) bool {
	ip := net.ParseIP(ipAddr)

	if ip == nil {
		return false
	}

	for _, block := range privateIPBlocks {
		if block.Contains(ip) {
			return true
		}
	}

	return false
}
