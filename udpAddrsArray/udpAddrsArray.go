package udpAddrsArray

import "net"

type UDPAddrsArray []*net.UDPAddr

func NewUDPAddrsArray() UDPAddrsArray {
	return make([]*net.UDPAddr, 0)
}

func (arr UDPAddrsArray) Filter(f func(addr *net.UDPAddr) bool) UDPAddrsArray {
	filteredArr := NewUDPAddrsArray()
	for _, v := range arr {
		if f(v) {
			filteredArr = append(filteredArr, v)
		}
	}
	return filteredArr
}

func (arr UDPAddrsArray) Add(addr *net.UDPAddr) UDPAddrsArray {
	return append(arr, addr)
}
