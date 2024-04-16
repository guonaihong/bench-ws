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
	return addrs[:n], nil
}

func GetFrameworkPidServerAddrs(framework string) (string, error) {
	ports, err := GetFrameworkBenchmarkPorts(framework)
	if err != nil {
		return "", err
	}
	addr := fmt.Sprintf(":%d", ports[len(ports)-1]+1)
	return addr, nil
}
