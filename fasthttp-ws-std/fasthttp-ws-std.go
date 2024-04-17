package main

import (
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/fasthttp/websocket"
	"github.com/guonaihong/bench-ws/config"
	"github.com/guonaihong/bench-ws/core"
	"github.com/guonaihong/clop"
)

type Config struct {
	// 打开性能优化开关
	UseReader      bool `clop:"short;long" usage:"use reader"`
	ReadBufferSize int  `clop:"short;long" usage:"read buffer size" default:"1024"`

	Addr string `clop:"short;long" usage:"websocket server address" default:":5555"`
	core.BaseCmd
}

var upgrader = websocket.Upgrader{}

func (c *Config) work(conn *websocket.Conn) {
	defer conn.Close()

	if !c.UseReader {
		for {
			mt, message, err := conn.ReadMessage()
			if err != nil {
				// log.Printf("read message failed: %v", err)
				return
			}
			err = conn.WriteMessage(mt, message)
			if err != nil {
				// log.Printf("write failed: %v", err)
				return
			}
		}
	}
	var nread int
	buffer := make([]byte, c.ReadBufferSize)
	readBuffer := buffer
	for {
		mt, reader, err := conn.NextReader()
		if err != nil {
			// log.Printf("read failed: %v", err)
			return
		}
		for {
			if nread == len(readBuffer) {
				readBuffer = append(readBuffer, buffer...)
			}
			n, err := reader.Read(readBuffer[nread:])
			nread += n
			if err == io.EOF {
				break
			}
		}
		err = conn.WriteMessage(mt, readBuffer[:nread])
		nread = 0
		if err != nil {
			// log.Printf("write failed: %v", err)
			return
		}
	}
}

func (c *Config) echo(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("upgrade failed: %v", err)
		return
	}
	conn.SetReadDeadline(time.Time{})

	go c.work(conn)
}

func (c *Config) startServers(addrs []string) []net.Listener {
	lns := make([]net.Listener, 0, len(addrs))
	for _, addr := range addrs {
		mux := &http.ServeMux{}
		mux.HandleFunc("/ws", c.echo)
		core.HandleCommon(mux)
		server := http.Server{
			// Addr:    addr,
			Handler: mux,
		}
		ln, err := core.Listen("tcp", addr, c.Reuse)
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
func main() {
	var cnf Config
	clop.Bind(&cnf)

	addrs, err := config.GetFrameworkServerAddrs(config.Fasthttp, cnf.LimitPortRange)
	if err != nil {
		log.Fatalf("GetFrameworkBenchmarkAddrs(%v) failed: %v", config.Fasthttp, err)
	}
	lns := core.StartServers(addrs, cnf.echo, cnf.Reuse)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	<-interrupt
	for _, ln := range lns {
		ln.Close()
	}
}
