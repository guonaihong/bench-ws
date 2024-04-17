package core

import (
	"net"

	"github.com/libp2p/go-reuseport"
)

func Listen(network, addr string, reuse bool) (net.Listener, error) {
	if reuse {
		return reuseport.Listen(network, addr)
	}
	return net.Listen(network, addr)
}

func Listen2(reuse bool) func(network, addr string) (net.Listener, error) {
	return func(network, addr string) (net.Listener, error) {
		if reuse {
			return reuseport.Listen(network, addr)
		}
		return net.Listen(network, addr)
	}
}
