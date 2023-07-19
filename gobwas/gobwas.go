package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/guonaihong/clop"
)

type Config struct{}

func (c *Config) work(conn net.Conn, brw *bufio.ReadWriter) {
	defer conn.Close()
	for {
		msg, op, err := wsutil.ReadClientData(brw)
		if err != nil {
			// log.Printf("read failed: %v", err)
			return
		}
		err = wsutil.WriteServerMessage(conn, op, msg)
		if err != nil {
			// log.Printf("write failed: %v", err)
			return
		}
	}
}

func (c *Config) echo(w http.ResponseWriter, r *http.Request) {
	conn, br, _, err := ws.UpgradeHTTP(r, w)
	if err != nil {
		log.Printf("upgrade failed: %v", err)
		return
	}
	conn.SetReadDeadline(time.Time{})

	go c.work(conn, br)
}

func main() {
	var conf Config
	clop.Bind(&conf)

	mux := &http.ServeMux{}
	mux.HandleFunc("/", conf.echo)

	go func() {
		// log.Println(http.ListenAndServe(":6060", nil))
	}()
	rawTCP, err := net.Listen("tcp", ":6001")
	if err != nil {
		fmt.Println("Listen fail:", err)
		return
	}

	log.Println("non-tls server exit:", http.Serve(rawTCP, mux))
}
