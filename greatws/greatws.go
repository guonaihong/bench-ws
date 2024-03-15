package main

import (
	_ "embed"
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/http"
	"sync/atomic"
	"time"

	_ "net/http/pprof"

	"github.com/antlabs/greatws"
	"github.com/guonaihong/clop"
)

type Config struct {
	RunInEventLoop bool   `clop:"short;long" usage:"run in event loop"`
	Addr           string `clop:"short;long" usage:"websocket server address" default:":9001"`
	EnableUtf8     bool   `clop:"short;long" usage:"enable utf8"`
	// 几倍的窗口大小
	WindowsMultipleTimesPayloadSize int `clop:"short;long" usage:"windows multiple times payload size"`
	// 打开tcp nodealy
	OpenTcpDelay bool `clop:"short;long" usage:"tcp delay"`
	// 使用stream模式, 一个连接对应一个go程
	StreamMode   bool   `clop:"short;long" usage:"use stream"`
	UnStreamMode bool   `clop:"short;long" usage:"use stream"`
	CustomMode   string `clop:"short;long" usage:"custom mode"`
	// 使用go程绑定模式, greatws默认模式
	GoRoutineBindMode bool `clop:"short;long" usage:"use go routine bind"`
	// 开启对流量压测友好的模式
	TrafficMode bool `clop:"short;long" usage:"enable pressure mode"`
	// 开启解析loop
	DisableParseLoop bool `clop:"short;long" usage:"disable parse loopo"`
	// 设置事件个数
	EventNum int `clop:"long" usage:"event number"`
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

var (
	total   uint64
	success uint64
)

func (e *echoHandler) OnMessage(c *greatws.Conn, op greatws.Opcode, msg []byte) {
	// fmt.Printf("OnMessage: %s, len(%d), op:%d\n", msg, len(msg), op)
	// if err := c.WriteTimeout(op, msg, 3*time.Second); err != nil {
	// 	fmt.Println("write fail:", err)
	// }
	atomic.AddUint64(&total, 1)
	if err := c.WriteMessage(op, msg); err != nil {
		slog.Error("write fail:", err)
	} else {
		atomic.AddUint64(&success, 1)
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

	evOpts := []greatws.EvOption{
		greatws.WithEventLoops(cnf.EventNum),
		greatws.WithBusinessGoNum(80, 10, 80),
		greatws.WithMaxEventNum(1000),
		greatws.WithLogLevel(slog.LevelError),
	}
	if cnf.TrafficMode {
		// evOpts = append(evOpts, greatws.WithBusinessGoTrafficMode())
	}
	if cnf.DisableParseLoop {
		evOpts = append(evOpts, greatws.WithDisableParseInParseLoop())
	}
	h.m = greatws.NewMultiEventLoopMust(evOpts...) // epoll, kqueue

	h.m.Start()

	opts := []greatws.ServerOption{
		greatws.WithServerReplyPing(),
		// greatws.WithServerDecompression(),
		greatws.WithServerIgnorePong(),
		greatws.WithServerCallback(&echoHandler{}),
		// greatws.WithServerEnableUTF8Check(),
		greatws.WithServerReadTimeout(5 * time.Second),
		greatws.WithServerMultiEventLoop(h.m),

		greatws.WithServerWindowsMultipleTimesPayloadSize(windowsSize),
	}

	switch {
	case cnf.RunInEventLoop:
		opts = append(opts, greatws.WithServerCallbackInEventLoop())
	case cnf.GoRoutineBindMode:
	case cnf.StreamMode:
		opts = append(opts, greatws.WithServerCustomTaskMode("stream"))
	case cnf.UnStreamMode:
		opts = append(opts, greatws.WithServerUnstreamMode())
	}

	if len(cnf.CustomMode) > 0 {
		opts = append(opts, greatws.WithServerCustomTaskMode(cnf.CustomMode))
	}

	upgrader = greatws.NewUpgrade(opts...)

	fmt.Printf("apiname:%s\n", h.m.GetApiName())

	go func() {
		for {
			time.Sleep(time.Second)
			fmt.Printf("OnMessage:%d, OnMessageSuccess:%d, curConn:%d, curGoNum:%d,curTask:%d, readSyscall:%d, writeSyscall:%d, realloc:%d, moveBytes:%d, readEv:%d writeEv:%d\n",
				atomic.LoadUint64(&total),
				atomic.LoadUint64(&success),
				h.m.GetCurConnNum(),
				h.m.GetCurGoNum(),
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
