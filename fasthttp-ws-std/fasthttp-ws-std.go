package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/fasthttp/websocket"
	"github.com/guonaihong/clop"
)

type Config struct {
	// 打开性能优化开关
	UseReader      bool `clop:"short;long" usage:"use reader"`
	ReadBufferSize int  `clop:"short;long" usage:"read buffer size" default:"1024"`

	Addr string `clop:"short;long" usage:"websocket server address" default:":5555"`
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

func main() {
	var conf Config
	clop.Bind(&conf)

	mux := &http.ServeMux{}
	mux.HandleFunc("/", conf.echo)

	go func() {
		// log.Println(http.ListenAndServe(":6060", nil))
	}()
	rawTCP, err := net.Listen("tcp", conf.Addr)
	if err != nil {
		fmt.Println("Listen fail:", err)
		return
	}

	log.Println("non-tls server exit:", http.Serve(rawTCP, mux))
}
