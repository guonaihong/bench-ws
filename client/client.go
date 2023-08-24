package main

import (
	"bytes"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	// _ "net/http/pprof"

	"github.com/antlabs/quickws"
	"github.com/guonaihong/clop"
)

// https://github.com/snapview/tokio-tungstenite/blob/master/examples/autobahn-client.rs

type Client struct {
	WSAddr string `clop:"short;long" usage:"websocket server address"`
	// 运行总次数
	Total int `clop:"short;long" usage:"total" default:"100"`

	PayloadSize int `clop:"short;long" usage:"payload size" default:"1024"`
	// 连接数
	Conns int `clop:"long" usage:"conns" default:"10000"`
	// 协程数
	Concurrency int `clop:"short;long" usage:"concurrency" default:"1000"`

	// 测试时间
	Duration time.Duration `clop:"short;long" usage:"duration"`

	OpenCheck bool `clop:"long" usage:"open check"`

	OpenTmpResult bool `clop:"long" usage:"open tmp result"`

	mu sync.Mutex

	result []int
}

var int64Count int64

var payload []byte

// var payload = []byte()

type echoHandler struct {
	// done chan struct{}
	data  chan struct{}
	total int
	curr  int

	OpenCheck bool
}

func (e *echoHandler) OnOpen(c *quickws.Conn) {
	c.WriteMessage(quickws.Binary, payload)
}

func (e *echoHandler) OnMessage(c *quickws.Conn, op quickws.Opcode, msg []byte) {
	// fmt.Println("OnMessage:", c, op, msg)
	if op == quickws.Text || op == quickws.Binary {
		c.WriteMessage(op, payload)
		if e.OpenCheck {
			if !bytes.Equal(msg, payload) {
				panic("payload not equal")
			}
		}

		atomic.AddInt64(&int64Count, 1)
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
	c, err := quickws.Dial(client.WSAddr,
		quickws.WithClientReplyPing(),
		// quickws.WithClientCompression(),
		// quickws.WithClientDecompressAndCompress(),
		quickws.WithClientCallback(&echoHandler{total: currTotal, data: data, OpenCheck: client.OpenCheck}),
		// quickws.WithClientCallback(&echoHandler{done: done}),
	)
	if err != nil {
		fmt.Println("Dial fail:", err)
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
		for i, v := range c.result {
			fmt.Printf("%ds:%d/s ", i+1, v)
		}
		fmt.Printf("\n")
	}()

	for i := 0; i < c.Concurrency; i++ {
		go func() {
			defer wg.Done()

			c.runTest(c.Total/c.Concurrency, data)
		}()
	}
}

func (c *Client) printTps(now time.Time, sec *int) {
	count := atomic.LoadInt64(&int64Count)
	n := int64(time.Since(now).Seconds())
	if n == 0 {
		n = 1
	}

	if c.OpenTmpResult {
		fmt.Printf("sec: %d, count: %d, tps: %d\n", *sec, count, count/n)
	}

	c.mu.Lock()
	c.result = append(c.result, int(count/n))
	c.mu.Unlock()
}

func (c *Client) Run(now time.Time) {
	for sec := 1; ; sec++ {
		time.Sleep(time.Second)
		c.printTps(now, &sec)
	}
}

func main() {
	var c Client

	go func() {
		// log.Println(http.ListenAndServe(":6060", nil))
	}()
	clop.MustBind(&c)
	payload = bytes.Repeat([]byte("𠜎"), c.PayloadSize/len("𠜎"))
	data := make(chan struct{}, c.Total)

	now := time.Now()
	go c.producer(data)
	go c.Run(now)
	c.consumer(data)
}
