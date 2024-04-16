package core

import (
	"go-websocket-benchmark/frameworks"
	"log"
	"net"
	"net/http"
)

func StartServers(addrs []string, echo http.HandlerFunc) []net.Listener {
	lns := make([]net.Listener, 0, len(addrs))
	for _, addr := range addrs {
		mux := &http.ServeMux{}
		mux.HandleFunc("/ws", echo)
		HandleCommon(mux)
		server := http.Server{
			// Addr:    addr,
			Handler: mux,
		}
		ln, err := frameworks.Listen("tcp", addr)
		if err != nil {
			log.Fatalf("Listen failed: %v", err)
		}
		lns = append(lns, ln)
		go func() {
			log.Printf("server exit: %v", server.Serve(ln))
		}()
	}
	return lns
}
