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
	"time"

	"github.com/fasthttp/websocket"
	"github.com/guonaihong/bench-ws/core"
	"github.com/guonaihong/bench-ws/pkg/port"
	"github.com/guonaihong/clop"
)

type Config struct {
	// 打开性能优化开关
	UseReader bool   `clop:"short;long" usage:"use reader"`
	Addr      string `clop:"short;long" usage:"websocket server address" default:":5555"`
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

	portRange, err := port.GetPortRange("FASTHTTP-WS-STD")
	if err != nil {
		log.Fatalf("GetPortRange(%v) failed: %v", "FASTHTTP-WS-STD", err)
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
