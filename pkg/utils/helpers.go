package utils

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/dicedb/dicedb-go/wire"
)

// ParseHostAndPort splits a URL string in format "host:port" and returns the host and port
func ParseHostAndPort(url string) (string, int) {
	// If the URL is not in the "host:port" format, treat
	// the URL as the host and use the default port 7379
	host := url
	port := 7379

	// If the URL contains a colon, try to split it into host and port
	if strings.Contains(url, ":") {
		var err error
		var portStr string

		host, portStr, err = net.SplitHostPort(url)
		if err == nil {
			portInt, err := strconv.Atoi(portStr)
			if err == nil {
				port = portInt
			}
		}
	}

	return host, port
}

// FormatDiceDBResponse formats the DiceDB response
func FormatDiceDBResponse(resp *wire.Response) string {
	if resp.Err != "" {
		return fmt.Sprintf("Error: %s", resp.Err)
	}

	// Handle different value types
	// TODO: Handle more cases like: https://github.com/DiceDB/dicedb-cli/blob/1c61ed7ec2a24f1483a59965df73450d575bbab6/ironhawk/main.go#L136
	switch resp.Value.(type) {
	case *wire.Response_VStr:
		return resp.GetVStr()
	case *wire.Response_VInt:
		return fmt.Sprintf("%d", resp.GetVInt())
	case *wire.Response_VFloat:
		return fmt.Sprintf("%f", resp.GetVFloat())
	case *wire.Response_VBytes:
		return fmt.Sprintf("%s", resp.GetVBytes())
	case *wire.Response_VNil:
		return "(nil)"
	default:
		return fmt.Sprintf("%v", resp)
	}
}
