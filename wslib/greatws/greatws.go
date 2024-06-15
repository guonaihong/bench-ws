package main

import (
	_ "embed"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"time"

	_ "net/http/pprof"

	"github.com/antlabs/greatws"
	"github.com/guonaihong/bench-ws/config"
	"github.com/guonaihong/bench-ws/core"
	"github.com/guonaihong/clop"
)

type Config struct {
	RunInEventLoop bool   `clop:"long" usage:"run in event loop"`
	Addr           string `clop:"short;long" usage:"websocket server address" default:":9001"`
	EnableUtf8     bool   `clop:"short;long" usage:"enable utf8"`
	// 几倍的窗口大小
	WindowsMultipleTimesPayloadSize int `clop:"short;long" usage:"windows multiple times payload size"`
	// 打开tcp nodealy
	OpenTcpDelay bool `clop:"short;long" usage:"tcp delay"`
	// 使用stream模式, 一个连接对应一个go程
	StreamMode bool   `clop:"short;long" usage:"use stream"`
	CustomMode string `clop:"short;long" usage:"custom mode"`
	// 使用go程绑定模式, greatws默认模式
	GoRoutineBindMode bool `clop:"short;long" usage:"use go routine bind"`
	// 开启对流量压测友好的模式
	TrafficMode bool `clop:"short;long" usage:"enable pressure mode"`
	// 开启解析loop
	DisableParseLoop bool `clop:"short;long" usage:"disable parse loopo"`
	// 设置事件个数
	EventNum int `clop:"long" usage:"event number"`
	MaxGoNum int `clop:"long" usage:"max go number" default:"80"`

	ProcessSleep time.Duration `clop:"long" usage:"process sleep"`
	core.BaseCmd
	m *greatws.MultiEventLoop
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
	if e.ProcessSleep > 0 {
		time.Sleep(e.ProcessSleep)
	}
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
	errMsg := "nil"
	if err != nil {
		errMsg = err.Error()
	}
	slog.Error("OnClose:", "err", errMsg)
}

func (h *Config) echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r)
	if err != nil {
		slog.Error("Upgrade fail:", "err", err.Error())
	}
	_ = c
}

func main() {

	var cnf Config
	clop.Bind(&cnf)

	go func() {
		log.Println(http.ListenAndServe(":6060", nil))
	}()

	windowsSize := float32(1.0)
	if cnf.WindowsMultipleTimesPayloadSize > 0 {
		windowsSize = float32(cnf.WindowsMultipleTimesPayloadSize)
	}

	initCount, minCount, maxCount := 80, 10, cnf.MaxGoNum
	evOpts := []greatws.EvOption{
		greatws.WithEventLoops(cnf.EventNum),
		greatws.WithBusinessGoNum(initCount, minCount, maxCount),
		greatws.WithMaxEventNum(1000),
		greatws.WithLogLevel(slog.LevelError),
	}

	fmt.Printf("init:%d, min:%d, max:%d\n", initCount, minCount, maxCount)
	if cnf.TrafficMode {
		// evOpts = append(evOpts, greatws.WithBusinessGoTrafficMode())
	}
	if cnf.DisableParseLoop {
		evOpts = append(evOpts, greatws.WithDisableParseInParseLoop())
	}
	cnf.m = greatws.NewMultiEventLoopMust(evOpts...) // epoll, kqueue

	cnf.m.Start()

	opts := []greatws.ServerOption{
		greatws.WithServerReplyPing(),
		// greatws.WithServerDecompression(),
		greatws.WithServerIgnorePong(),
		greatws.WithServerCallback(&echoHandler{Config: &cnf}),
		// greatws.WithServerEnableUTF8Check(),
		// greatws.WithServerReadTimeout(60 * time.Second),
		greatws.WithServerMultiEventLoop(cnf.m),

		greatws.WithServerWindowsMultipleTimesPayloadSize(windowsSize),
	}

	switch {
	case cnf.RunInEventLoop:
		opts = append(opts, greatws.WithServerCallbackInEventLoop())
	case cnf.GoRoutineBindMode:
	case cnf.StreamMode:
		opts = append(opts, greatws.WithServerCustomTaskMode("stream"))
	}

	if len(cnf.CustomMode) > 0 {
		opts = append(opts, greatws.WithServerCustomTaskMode(cnf.CustomMode))
	}

	upgrader = greatws.NewUpgrade(opts...)

	fmt.Printf("apiname:%s\n", cnf.m.GetApiName())

	go func() {
		for {
			time.Sleep(time.Second)
			fmt.Printf("OnMessage:%d, OnMessageSuccess:%d, curConn:%d, curGoNum:%d,curTask:%d, readSyscall:%d, writeSyscall:%d, realloc:%d, moveBytes:%d, readEv:%d writeEv:%d\n",
				atomic.LoadUint64(&total),
				atomic.LoadUint64(&success),
				cnf.m.GetCurConnNum(),
				cnf.m.GetCurGoNum(),
				cnf.m.GetCurTaskNum(),
				cnf.m.GetReadSyscallNum(),
				cnf.m.GetWriteSyscallNum(),
				cnf.m.GetReallocNum(),
				cnf.m.GetMoveBytesNum(),
				cnf.m.GetReadEvNum(),
				cnf.m.GetWriteEvNum(),
			)
		}
	}()

	addrs, err := config.GetFrameworkServerAddrs(config.Greatws, cnf.LimitPortRange)
	if err != nil {
		log.Fatalf("GetFrameworkBenchmarkAddrs(%v) failed: %v", config.Quickws, err)
	}

	lns := core.StartServers(addrs, cnf.echo, cnf.Reuse)
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	<-interrupt
	for _, ln := range lns {
		ln.Close()
	}
}
