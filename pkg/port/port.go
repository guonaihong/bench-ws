package port

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// PortRange represents a range of ports
type PortRange struct {
	Start int
	End   int
}

// Default port ranges for different libraries
var defaultPortRanges = map[string]PortRange{
	"QUICKWS":          {Start: 2000, End: 2000},
	"GREATWS":          {Start: 2100, End: 2100},
	"NHOOYR":           {Start: 2200, End: 2200},
	"GWS":              {Start: 2300, End: 2300},
	"GWS-STD":          {Start: 2400, End: 2400},
	"GORILLA":          {Start: 2500, End: 2500},
	"GOBWAS":           {Start: 2600, End: 2600},
	"NETTYWS":          {Start: 2700, End: 2700},
	"NBIO-STD":         {Start: 2800, End: 2800},
	"NBIO-NONBLOCKING": {Start: 2900, End: 2900},
	"NBIO-BLOCKING":    {Start: 3000, End: 3000},
	"NBIO-MIXED":       {Start: 3100, End: 3100},
	"FASTHTTP-WS-STD":  {Start: 3200, End: 3200},
	"HERTZ-STD":        {Start: 3300, End: 3300},
	"HERTZ":            {Start: 3400, End: 3400},
	"GREATWS-EVENT":    {Start: 3500, End: 3500},
}

// GetPortRange returns the port range for a given library name
func GetPortRange(libName string) (*PortRange, error) {
	libName = strings.ToUpper(libName)
	startPort := os.Getenv(fmt.Sprintf("%s_START_PORT", libName))
	endPort := os.Getenv(fmt.Sprintf("%s_END_PORT", libName))

	// If environment variables are not set, use default values
	if startPort == "" || endPort == "" {
		if defaultRange, ok := defaultPortRanges[libName]; ok {
			return &defaultRange, nil
		}
		return nil, fmt.Errorf("port range not configured for library %s", libName)
	}

	start, err := strconv.Atoi(startPort)
	if err != nil {
		return nil, fmt.Errorf("invalid start port for %s: %v", libName, err)
	}

	end, err := strconv.Atoi(endPort)
	if err != nil {
		return nil, fmt.Errorf("invalid end port for %s: %v", libName, err)
	}

	return &PortRange{
		Start: start,
		End:   end,
	}, nil
}

// GetNextPort returns the next available port in the range
func (pr *PortRange) GetNextPort(current int) int {
	if current < pr.Start || current >= pr.End {
		return pr.Start
	}
	return current + 1
}

// Format returns the port range as a string
func (pr *PortRange) Format() string {
	return fmt.Sprintf("%d-%d", pr.Start, pr.End)
}
