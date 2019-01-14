package udpAddrsArray

import (
	"bytes"
	"net"
)

type UDPAddrsArray []*net.UDPAddr

func NewUDPAddrsArray() UDPAddrsArray {
	return make([]*net.UDPAddr, 0)
}

func (arr UDPAddrsArray) Remove(addr *net.UDPAddr) UDPAddrsArray {
	filteredArr := NewUDPAddrsArray()
	for _, addrI := range arr {
		if bytes.Compare(addrI.IP, addr.IP) != 0 {
			filteredArr = append(filteredArr, addrI)
		}
	}
	return filteredArr
}

func (arr UDPAddrsArray) Add(addr *net.UDPAddr) UDPAddrsArray {
	return append(arr, addr)
}
