package main

import (
	"bufio"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/guonaihong/bench-ws/config"
	"github.com/guonaihong/bench-ws/core"
	"github.com/guonaihong/clop"
)

type Config struct {
	Addr string `clop:"short;long" usage:"websocket server address" default:":9001"`
	// 使用限制端口范围, 默认1， -1表示不限制
	LimitPortRange int `clop:"short;long" usage:"limit port range" default:"1"`
}

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
	var cnf Config
	clop.Bind(&cnf)

	mux := &http.ServeMux{}
	mux.HandleFunc("/", cnf.echo)

	addrs, err := config.GetFrameworkServerAddrs(config.Gobwas, cnf.LimitPortRange)
	if err != nil {
		log.Fatalf("GetFrameworkBenchmarkAddrs(%v) failed: %v", config.Gobwas, err)
	}
	lns := core.StartServers(addrs, cnf.echo)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	<-interrupt
	for _, ln := range lns {
		ln.Close()
	}
}
