package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/coder/websocket"
	"github.com/guonaihong/bench-ws/pkg/port"
)

func main() {

	portRange, err := port.GetPortRange("NHOOYR")
	if err != nil {
		log.Fatalf("GetPortRange(%v) failed: %v", "NHOOYR", err)
	}

	fmt.Println("NHOOYR server starting on ports", portRange.Start, "-", portRange.End)

	wg := sync.WaitGroup{}
	defer wg.Wait()
	for port := portRange.Start; port <= portRange.End; port++ {
		wg.Add(1)
		go startServer(port, &wg)
	}

	// Wait a moment for server to start
	time.Sleep(time.Second)

}

// startServer starts a WebSocket server on port 8080
func startServer(port int, wg *sync.WaitGroup) {
	defer wg.Done()

	mux := http.NewServeMux()
	mux.HandleFunc("/", handleWebSocket)

	fmt.Println("NHOOYR server starting on :", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), mux))
}

func echo(c *websocket.Conn) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	typ, r, err := c.Reader(ctx)
	if err != nil {
		return err
	}

	w, err := c.Writer(ctx, typ)
	if err != nil {
		return err
	}

	_, err = io.Copy(w, r)
	if err != nil {
		return fmt.Errorf("failed to io.Copy: %w", err)
	}

	err = w.Close()
	return err
}

// handleWebSocket handles WebSocket connections
func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Accept the WebSocket connection
	c, err := websocket.Accept(w, r, nil)
	if err != nil {
		log.Printf("Failed to accept WebSocket connection: %v", err)
		return
	}
	defer c.CloseNow()

	fmt.Println("New WebSocket connection established")

	for {
		err = echo(c)
		if websocket.CloseStatus(err) == websocket.StatusNormalClosure {
			return
		}
		if err != nil {
			fmt.Printf("failed to echo with %v: %v", r.RemoteAddr, err)
			return
		}
	}

	// Close connection gracefully
	c.Close(websocket.StatusNormalClosure, "Server closing connection")
}
