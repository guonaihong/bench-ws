package main

import (
	_ "embed"
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/http"
	"runtime"
	"time"

	_ "net/http/pprof"

	"github.com/antlabs/greatws"
	"github.com/guonaihong/clop"
)

type Config struct {
	Addr string `clop:"short;long" usage:"websocket server address" default:":9001"`

	EnableUtf8 bool `clop:"short;long" usage:"enable utf8"`
	// 几倍的窗口大小
	WindowsMultipleTimesPayloadSize int `clop:"short;long" usage:"windows multiple times payload size"`
	// 打开tcp nodealy
	OpenTcpDelay bool `clop:"short;long" usage:"tcp delay"`
}

var upgrader *greatws.UpgradeServer

// go:embed public.crt
// var certPEMBlock []byte

// go:embed privatekey.pem
// var keyPEMBlock []byte
type echoHandler struct {
	*Config
}

func (e *echoHandler) OnOpen(c *greatws.Conn) {
	// fmt.Printf("OnOpen: %p\n", c)
}

func (e *echoHandler) OnMessage(c *greatws.Conn, op greatws.Opcode, msg []byte) {
	// fmt.Printf("OnMessage: %s, len(%d), op:%d\n", msg, len(msg), op)
	// if err := c.WriteTimeout(op, msg, 3*time.Second); err != nil {
	// 	fmt.Println("write fail:", err)
	// }
	if err := c.WriteMessage(op, msg); err != nil {
		slog.Error("write fail:", err)
	}
}

func (e *echoHandler) OnClose(c *greatws.Conn, err error) {
	// errMsg := ""
	// if err != nil {
	// 	errMsg = err.Error()
	// }
	// slog.Error("OnClose:", errMsg)
}

type handler struct {
	m *greatws.MultiEventLoop
}

func (h *handler) echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r)
	if err != nil {
		slog.Error("Upgrade fail:", "err", err.Error())
	}
	_ = c
}

func main() {
	var h handler

	var cnf Config
	clop.Bind(&cnf)

	go func() {
		log.Println(http.ListenAndServe(":6060", nil))
	}()

	windowsSize := float32(1.0)
	if cnf.WindowsMultipleTimesPayloadSize > 0 {
		windowsSize = float32(cnf.WindowsMultipleTimesPayloadSize)
	}
	// debug io-uring
	// h.m = greatws.NewMultiEventLoopMust(greatws.WithEventLoops(0), greatws.WithMaxEventNum(1000), greatws.WithIoUring(), greatws.WithLogLevel(slog.LevelDebug))
	h.m = greatws.NewMultiEventLoopMust(
		greatws.WithEventLoops(runtime.NumCPU()*2),
		greatws.WithBusinessGoNum(80, 10, 10000),
		greatws.WithMaxEventNum(1000),
		greatws.WithLogLevel(slog.LevelError)) // epoll, kqueue
	h.m.Start()

	upgrader = greatws.NewUpgrade(
		greatws.WithServerReplyPing(),
		// greatws.WithServerDecompression(),
		greatws.WithServerIgnorePong(),
		greatws.WithServerCallback(&echoHandler{}),
		// greatws.WithServerEnableUTF8Check(),
		greatws.WithServerReadTimeout(5*time.Second),
		greatws.WithServerMultiEventLoop(h.m),

		greatws.WithServerWindowsMultipleTimesPayloadSize(windowsSize),
	)

	fmt.Printf("apiname:%s\n", h.m.GetApiName())

	go func() {
		for {
			time.Sleep(time.Second)
			fmt.Printf("curConn:%d, curTask:%d, readSyscall:%d, writeSyscall:%d, realloc:%d, moveBytes:%d, readEv:%d writeEv:%d\n",
				h.m.GetCurConnNum(),
				h.m.GetCurTaskNum(),
				h.m.GetReadSyscallNum(),
				h.m.GetWriteSyscallNum(),
				h.m.GetReallocNum(),
				h.m.GetMoveBytesNum(),
				h.m.GetReadEvNum(),
				h.m.GetWriteEvNum(),
			)
		}
	}()

	mux := &http.ServeMux{}
	mux.HandleFunc("/", h.echo)

	rawTCP, err := net.Listen("tcp", cnf.Addr)
	if err != nil {
		fmt.Println("Listen fail:", err)
		return
	}

	log.Println("non-tls server exit:", http.Serve(rawTCP, mux))

	// cert, err := tls.X509KeyPair(certPEMBlock, keyPEMBlock)
	// if err != nil {
	// 	log.Fatalf("tls.X509KeyPair failed: %v", err)
	// }
	// tlsConfig := &tls.Config{
	// 	Certificates:       []tls.Certificate{cert},
	// 	InsecureSkipVerify: true,
	// }
	// lnTLS, err := tls.Listen("tcp", "localhost:9002", tlsConfig)
	// if err != nil {
	// 	panic(err)
	// }
	// log.Println("tls server exit:", http.Serve(lnTLS, mux))
}
