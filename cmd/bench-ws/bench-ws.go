package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"
	"sync/atomic"
	"time"

	_ "net/http/pprof"

	"github.com/antlabs/quickws"
	"github.com/guonaihong/bench-ws/config"
	"github.com/guonaihong/bench-ws/report"
	"github.com/guonaihong/clop"
)

// https://github.com/snapview/tokio-tungstenite/blob/master/examples/autobahn-client.rs

type Client struct {
	WSAddr         string        `clop:"short;long" usage:"WebSocket server address (e.g., ws://host:port or ws://host:minport-maxport)" default:""`
	Name           string        `clop:"short;long" usage:"Server name" default:""`
	Label          string        `clop:"long" usage:"Title of the chart for the line graph" default:""`
	Total          int           `clop:"short;long" usage:"Total number of runs" default:"100"`
	PayloadSize    int           `clop:"short;long" usage:"Size of the payload" default:"1024"`
	Conns          int           `clop:"long" usage:"Number of connections" default:"10000"`
	Concurrency    int           `clop:"short;long" usage:"Number of concurrent goroutines" default:"1000"`
	Duration       time.Duration `clop:"short;long" usage:"Duration of the test"`
	OpenCheck      bool          `clop:"long" usage:"Perform open check"`
	OpenTmpResult  bool          `clop:"long" usage:"Display temporary result"`
	JSON           bool          `clop:"long" usage:"Output JSON result"`
	Text           string        `clop:"long" usage:"Text to send"`
	SaveErr        bool          `clop:"long" usage:"Save error log"`
	LimitPortRange int           `clop:"short;long" usage:"Limit port range (1 for limited, -1 for unlimited)" default:"1"`
	mu             sync.Mutex

	result []int

	addrs []string
	index int64

	ctx    context.Context
	cancel context.CancelCauseFunc
}

func (c *Client) getAddrs() string {
	curIndex := int(atomic.AddInt64(&c.index, 1))
	return c.addrs[curIndex%len(c.addrs)]
}

var recvCount int64
var sendCount int64

var payload []byte

// var payload = []byte()

type echoHandler struct {
	// done chan struct{}
	data  chan struct{}
	total int
	curr  int

	*Client
}

// OnOpen is a callback function that is called when a WebSocket connection
// is established. It sends a binary message containing the payload to the
// server. It also increments the `sendCount` atomic integer.
//
// Parameters:
//   - c: A pointer to a `quickws.Conn` object representing the WebSocket
//     connection.
func (e *echoHandler) OnOpen(c *quickws.Conn) {
	// Send a binary message containing the payload to the server.
	c.WriteMessage(quickws.Binary, payload)
	// Increment the `sendCount` atomic integer.
	atomic.AddInt64(&sendCount, 1)
}

func (e *echoHandler) OnMessage(c *quickws.Conn, op quickws.Opcode, msg []byte) {
	atomic.AddInt64(&sendCount, 1)
	if op == quickws.Text || op == quickws.Binary {
		c.WriteMessage(op, payload)
		if e.OpenCheck {
			if !bytes.Equal(msg, payload) {
				if e.SaveErr {

					os.WriteFile(fmt.Sprintf("%x.err.log", c), payload, 0644)
					os.WriteFile(fmt.Sprintf("%v.success.log", c), msg, 0644)
				}
				panic("payload not equal")
			}
		}

		atomic.AddInt64(&recvCount, 1)
		select {
		case _, ok := <-e.data:
			if !ok {
				c.Close()
				return
			}
		default:
		}
	}
}

func (e *echoHandler) OnClose(c *quickws.Conn, err error) {
	// close(e.done)
}

func (client *Client) runTest(currTotal int, data chan struct{}) {
	curAddr := client.getAddrs()
	c, err := quickws.Dial(client.getAddrs(),
		quickws.WithClientReplyPing(),
		// quickws.WithClientCompression(),
		// quickws.WithClientDecompressAndCompress(),
		quickws.WithClientCallback(&echoHandler{total: currTotal, data: data, Client: client}),
		// quickws.WithClientCallback(&echoHandler{done: done}),
	)
	if err != nil {
		fmt.Printf("Dial %s, fail:%v\n", curAddr, err)
		return
	}

	c.ReadLoop()
}

// 生产者
func (c *Client) producer(data chan struct{}) {
	defer func() {
		close(data)

		if c.OpenTmpResult {
			fmt.Printf("bye bye producer")
		}
	}()
	if c.Duration > 0 {
		tk := time.NewTicker(c.Duration)
		for {
			select {
			case <-tk.C:
				// 时间到了
				// 排空chan
				for {
					select {
					case <-data:
					default:
						return
					}
				}
			case data <- struct{}{}:
			}
		}
	} else {
		for i := 0; i < c.Total; i++ {
			data <- struct{}{}
		}
	}
}

// 消费者
func (c *Client) consumer(data chan struct{}) {
	var wg sync.WaitGroup
	wg.Add(c.Concurrency)
	defer func() {
		wg.Wait()
		c.cancel(errors.New("wait all consumer done"))
		if !c.JSON {
			for i, v := range c.result {
				fmt.Printf("%ds:%d/s ", i+1, v)
			}
		}
		fmt.Printf("\n")
	}()

	for i := 0; i < c.Concurrency; i++ {
		go func(i int) {
			defer wg.Done()

			c.runTest(c.Total/c.Concurrency, data)
			// fmt.Printf("bye bye consumer:%d\n", i)
		}(i)
	}
}

func (c *Client) printTps(now time.Time, sec *int) {
	recvCount := atomic.LoadInt64(&recvCount)
	sendCount := atomic.LoadInt64(&sendCount)
	n := int64(time.Since(now).Seconds())
	if n == 0 {
		n = 1
	}

	if c.OpenTmpResult {
		fmt.Printf("sec: %d, recv-count: %d, send-count:%d recv-tps: %d, send-tps: %d\n", *sec, recvCount, sendCount, recvCount/n, sendCount/n)
	}

	c.mu.Lock()
	c.result = append(c.result, int(recvCount/n))
	c.mu.Unlock()
}

func (c *Client) Run(now time.Time) {
	nw := time.NewTicker(time.Second)
	sec := 1
	for {
		select {
		case <-nw.C:
			c.printTps(now, &sec)
			sec++
			nw.Reset(time.Second)
		case <-c.ctx.Done():
			if c.JSON {
				var d report.Dataset
				d.Label = c.Label
				d.Data = c.result
				d.Tension = 0.1
				all, err := json.Marshal(d)
				if err != nil {
					panic(err)
				}

				os.Stdout.Write(all)
			}
			return
		}
	}
}

func main() {
	var c Client

	clop.MustBind(&c)
	if len(c.Text) > 0 {
		payload = []byte(c.Text)
	} else {
		payload = bytes.Repeat([]byte("𠜎"), c.PayloadSize/len("𠜎"))
	}

	if c.WSAddr == "" && c.Name == "" {
		fmt.Printf("wsaddr or name is required, ./bench-ws -h\n")
		os.Exit(1)
	}
	c.addrs = config.GenerateAddrs(c.WSAddr, c.Name)
	data := make(chan struct{}, c.Total)

	now := time.Now()
	c.ctx, c.cancel = context.WithCancelCause(context.Background())
	var wg sync.WaitGroup
	wg.Add(2)
	defer wg.Wait()
	go func() {
		defer wg.Done()
		c.producer(data)
	}()
	go func() {
		defer wg.Done()
		c.Run(now)
	}()
	c.consumer(data)
}
