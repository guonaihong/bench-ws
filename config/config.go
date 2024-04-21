package config

// copy from go-websocket-benchmark
import (
	"fmt"
	"strconv"
	"strings"
)

const (
	Fasthttp           = "fasthttp"
	Gobwas             = "gobwas"
	Gorilla            = "gorilla"
	Gws                = "gws"
	GwsStd             = "gws_std"
	Hertz              = "hertz"
	HertzStd           = "hertz_std"
	NbioModBlocking    = "nbio_blocking"
	NbioModMixed       = "nbio_mixed"
	NbioModNonblocking = "nbio_nonblocking"
	NbioStd            = "nbio_std"
	GoNettyWs          = "nettyws"
	Nhooyr             = "nhooyr"
	Quickws            = "quickws"
	Greatws            = "greatws"
	GreatwsEvent       = "greatws_event"
)

var Ports = map[string]string{
	Fasthttp:           "10001:10050",
	Gobwas:             "11001:11050",
	Gorilla:            "12001:12050",
	Gws:                "13001:13050",
	GwsStd:             "14001:14050",
	Hertz:              "15001:15050",
	HertzStd:           "16001:16050",
	NbioModBlocking:    "17001:17050",
	NbioModMixed:       "18001:18050",
	NbioModNonblocking: "19001:19050",
	NbioStd:            "20001:20050",
	GoNettyWs:          "21001:21050",
	Nhooyr:             "22001:22050",
	Quickws:            "23001:23050",
	Greatws:            "24001:24050",
	GreatwsEvent:       "25001:25050",
}

var FrameworkList = []string{
	Fasthttp,
	Gobwas,
	Gorilla,
	Gws,
	GwsStd,
	Hertz,
	HertzStd,
	NbioModBlocking,
	NbioModMixed,
	NbioModNonblocking,
	NbioStd,
	GoNettyWs,
	Nhooyr,
	Quickws,
	Greatws,
	GreatwsEvent,
}

func GenerateAddrs(WSAddr, Name string) []string {
	var Addrs []string

	// 如果 WSAddr 为空并且 Name 不为空，将 host 设置为 "127.0.0.1"
	host := "ws://127.0.0.1"
	portAndPath := ""
	var parts []string
	if WSAddr != "" {
		if strings.HasPrefix(WSAddr, "ws://") {
			parts = strings.Split(WSAddr[5:], ":")
			host = fmt.Sprintf("ws://%s", parts[0])
		} else if strings.HasPrefix(WSAddr, "ws://") {
			parts := strings.Split(WSAddr[6:], ":")
			host = fmt.Sprintf("wss://%s", parts[0])
		} else {
			parts := strings.Split(WSAddr, ":")
			host = parts[0]
		}

		if len(parts) > 1 {
			portAndPath = parts[1]
		}
	}

	if portAndPath != "" {
		path := ""
		if strings.Contains(portAndPath, "/") {
			pos := strings.Index(portAndPath, "/")
			path = portAndPath[pos+1:]
			portAndPath = portAndPath[:pos]
		}

		if strings.Contains(portAndPath, "-") {
			// WSAddr 格式是 host:minport-maxport
			rangeParts := strings.Split(portAndPath, "-")
			minPortStr := rangeParts[0]
			maxPortStr := rangeParts[1]
			minPort, _ := strconv.Atoi(minPortStr)
			maxPort, _ := strconv.Atoi(maxPortStr)
			for i := minPort; i <= maxPort; i++ {
				if len(path) == 0 {
					Addrs = append(Addrs, fmt.Sprintf("%s:%d", host, i))
				} else {
					Addrs = append(Addrs, fmt.Sprintf("%s:%d/%s", host, i, path))
				}
			}
		} else {
			// WSAddr 格式是 host:port
			Addrs = append(Addrs, WSAddr)
		}
	} else if Name != "" {
		portRange, exists := Ports[Name]
		if exists {
			minPortStr := strings.Split(portRange, ":")[0]
			maxPortStr := strings.Split(portRange, ":")[1]
			minPort, _ := strconv.Atoi(minPortStr)
			maxPort, _ := strconv.Atoi(maxPortStr)
			for i := minPort; i <= maxPort; i++ {
				Addrs = append(Addrs, fmt.Sprintf("%s:%d", host, i))
			}
		}
	}

	if len(Addrs) == 0 {
		return []string{WSAddr}
	}
	return Addrs
}

func GetFrameworkBenchmarkPorts(framework string) ([]int, error) {
	portRange := strings.Split(Ports[framework], ":")
	minPort, err := strconv.Atoi(portRange[0])
	if err != nil {
		return nil, err
	}
	maxPort, err := strconv.Atoi(portRange[1])
	if err != nil {
		return nil, err
	}
	ports := []int{}
	for i := minPort; i <= maxPort; i++ {
		ports = append(ports, i)
	}
	return ports, nil
}

// GetFrameworkClientAddrs returns the client addresses for the given framework.
func GetFrameworkClientAddrs(framework string, host string, limit int) ([]string, error) {
	ports, err := GetFrameworkBenchmarkPorts(framework)
	if err != nil {
		return nil, err
	}
	addrs := make([]string, 0, len(ports))
	for _, port := range ports {
		addrs = append(addrs, fmt.Sprintf("%s:%d", host, port))
	}

	n := len(addrs)
	if limit > 0 {
		n = min(n, limit)
	}
	return addrs[:n], nil
}

func GetFrameworkServerAddrs(framework string, limit int) ([]string, error) {
	ports, err := GetFrameworkBenchmarkPorts(framework)
	if err != nil {
		return nil, err
	}
	addrs := make([]string, 0, len(ports))
	for _, port := range ports {
		addrs = append(addrs, fmt.Sprintf(":%d", port))
	}

	n := len(addrs)
	if limit > 0 {
		n = min(n, limit)
	}
	// framework的端口范围是
	addrs = addrs[:n]
	fmt.Printf("%s, %s-%s\n", framework, addrs[0], addrs[len(addrs)-1])
	return addrs, nil
}

func GetFrameworkPidServerAddrs(framework string) (string, error) {
	ports, err := GetFrameworkBenchmarkPorts(framework)
	if err != nil {
		return "", err
	}
	addr := fmt.Sprintf(":%d", ports[len(ports)-1]+1)
	return addr, nil
}
