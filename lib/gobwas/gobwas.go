package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/guonaihong/bench-ws/core"
	"github.com/guonaihong/bench-ws/pkg/port"
	"github.com/guonaihong/clop"
)

type Config struct {
	// 打开性能优化开关
	UseReader bool   `clop:"short;long" usage:"use reader"`
	Addr      string `clop:"short;long" usage:"websocket server address" default:":5555"`
	// 打开tcp nodealy
	OpenTcpDelay bool `clop:"short;long" usage:"tcp delay"`
	core.BaseCmd
}

func setNoDelay(c net.Conn, noDelay bool) error {
	if tcp, ok := c.(*net.TCPConn); ok {
		return tcp.SetNoDelay(noDelay)
	}
	return nil
}

func (c *Config) work(conn net.Conn) {
	defer conn.Close()

	if !c.UseReader {
		for {
			msg, op, err := wsutil.ReadClientData(conn)
			if err != nil {
				// log.Printf("read message failed: %v", err)
				return
			}
			err = wsutil.WriteServerMessage(conn, op, msg)
			if err != nil {
				// log.Printf("write failed: %v", err)
				return
			}
		}
	}

	reader := wsutil.NewReader(conn, ws.StateServerSide)
	writer := wsutil.NewWriter(conn, ws.StateServerSide, ws.OpText)
	buffer := make([]byte, c.ReadBufferSize)

	for {
		header, err := reader.NextFrame()
		if err != nil {
			// log.Printf("read header failed: %v", err)
			return
		}

		writer.Reset(conn, ws.StateServerSide, header.OpCode)
		if _, err = io.CopyBuffer(writer, reader, buffer); err != nil {
			// log.Printf("copy failed: %v", err)
			return
		}
		if err = writer.Flush(); err != nil {
			// log.Printf("flush failed: %v", err)
			return
		}
	}
}

func (c *Config) echo(w http.ResponseWriter, r *http.Request) {
	conn, _, _, err := ws.UpgradeHTTP(r, w)
	if err != nil {
		log.Printf("upgrade failed: %v", err)
		return
	}

	setNoDelay(conn, !c.OpenTcpDelay)
	go c.work(conn)
}

func (c *Config) startServer(port int, wg *sync.WaitGroup) {
	defer wg.Done()

	mux := &http.ServeMux{}
	mux.HandleFunc("/", c.echo)

	server := http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("Listen failed: %v", err)
	}

	log.Printf("server exit: %v", server.Serve(ln))
}

func main() {
	var cnf Config
	clop.Bind(&cnf)

	portRange, err := port.GetPortRange("GOBWAS")
	if err != nil {
		log.Fatalf("GetPortRange(%v) failed: %v", "GOBWAS", err)
	}
	fmt.Println("GOBWAS server starting on ports", portRange.Start, "-", portRange.End)

	wg := sync.WaitGroup{}
	for port := portRange.Start; port <= portRange.End; port++ {
		wg.Add(1)
		go cnf.startServer(port, &wg)
	}
	wg.Wait()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	<-interrupt
}
