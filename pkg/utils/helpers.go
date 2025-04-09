package utils

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/dicedb/dicedb-go/wire"
	"google.golang.org/protobuf/types/known/structpb"
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

	var result strings.Builder

	// Copied from: https://github.com/DiceDB/dicedb-cli/blob/1c61ed7ec2a24f1483a59965df73450d575bbab6/ironhawk/main.go#L136
	// Handle attributes if present
	if len(resp.Attrs.AsMap()) > 0 {
		attrs := []string{}
		for k, v := range resp.Attrs.AsMap() {
			attrs = append(attrs, fmt.Sprintf("%s=%s", k, v))
		}
		result.WriteString(fmt.Sprintf("[%s] ", strings.Join(attrs, ", ")))
	}

	// Handle string-string map if present
	if len(resp.VSsMap) > 0 {
		if result.Len() > 0 {
			result.WriteString("\n")
		}
		for k, v := range resp.VSsMap {
			result.WriteString(fmt.Sprintf("%s=%s\n", k, v))
		}
	}

	// Handle the primary value based on its type
	switch resp.Value.(type) {
	case *wire.Response_VStr:
		result.WriteString(resp.GetVStr())
	case *wire.Response_VInt:
		result.WriteString(fmt.Sprintf("%d", resp.GetVInt()))
	case *wire.Response_VFloat:
		result.WriteString(fmt.Sprintf("%f", resp.GetVFloat()))
	case *wire.Response_VBytes:
		result.WriteString(fmt.Sprintf("%s", resp.GetVBytes()))
	case *wire.Response_VNil:
		result.WriteString("(nil)")
	}

	// Handle list values if present
	if len(resp.GetVList()) > 0 {
		if result.Len() > 0 {
			result.WriteString("\n")
		}

		for i, v := range resp.GetVList() {
			switch v.GetKind().(type) {
			case *structpb.Value_NullValue:
				result.WriteString(fmt.Sprintf("%d) (nil)\n", i+1))
			case *structpb.Value_NumberValue:
				result.WriteString(fmt.Sprintf("%d) %f\n", i+1, v.GetNumberValue()))
			case *structpb.Value_StringValue:
				result.WriteString(fmt.Sprintf("%d) \"%s\"\n", i+1, v.GetStringValue()))
			case *structpb.Value_BoolValue:
				result.WriteString(fmt.Sprintf("%d) %t\n", i+1, v.GetBoolValue()))
			}
		}
	}

	return strings.TrimSpace(result.String())
}
