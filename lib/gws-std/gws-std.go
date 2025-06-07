package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/guonaihong/bench-ws/core"
	"github.com/guonaihong/bench-ws/pkg/port"
	"github.com/guonaihong/clop"
	"github.com/lxzan/gws"
)

type Config struct {
	// 打开性能优化开关
	UseReader bool   `clop:"short;long" usage:"use reader"`
	Addr      string `clop:"short;long" usage:"websocket server address" default:":5555"`
	// 打开tcp nodealy
	OpenTcpDelay bool `clop:"short;long" usage:"tcp delay"`
	core.BaseCmd
}

type Handler struct {
	*Config
}

func (h *Handler) OnOpen(socket *gws.Conn) {
	socket.SetReadDeadline(time.Time{})
}

func (h *Handler) OnClose(socket *gws.Conn, err error) {
}

func (h *Handler) OnPing(socket *gws.Conn, payload []byte) {
	socket.WritePong(payload)
}

func (h *Handler) OnPong(socket *gws.Conn, payload []byte) {
}

func (h *Handler) OnMessage(socket *gws.Conn, message *gws.Message) {
	socket.WriteMessage(message.Opcode, message.Bytes())
}

func setNoDelay(c net.Conn, noDelay bool) error {
	if tcp, ok := c.(*net.TCPConn); ok {
		return tcp.SetNoDelay(noDelay)
	}
	return nil
}

func (c *Config) echo(w http.ResponseWriter, r *http.Request) {
	upgrader := gws.NewUpgrader(&Handler{Config: c}, &gws.ServerOption{
		ReadBufferSize:  c.ReadBufferSize,
		WriteBufferSize: c.ReadBufferSize,
	})

	conn, err := upgrader.Upgrade(w, r)
	if err != nil {
		log.Printf("upgrade failed: %v", err)
		return
	}

	setNoDelay(conn.NetConn(), !c.OpenTcpDelay)
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

	go func() {
		log.Printf("server exit: %v", server.Serve(ln))
	}()
}

func main() {
	var cnf Config
	clop.Bind(&cnf)

	portRange, err := port.GetPortRange("GWS-STD")
	if err != nil {
		log.Fatalf("GetPortRange(%v) failed: %v", "GWS-STD", err)
	}

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
