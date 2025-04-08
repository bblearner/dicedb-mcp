package main

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/dicedb/dicedb-go"
	"github.com/dicedb/dicedb-go/wire"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	// Create a new MCP server
	s := server.NewMCPServer(
		"DiceDB MCP",
		"0.1.0",
		server.WithResourceCapabilities(true, true),
		server.WithLogging(),
	)

	// Create a new pingTool to ping the DiceDB server
	pingTool := mcp.NewTool("ping",
		mcp.WithDescription("Ping the DiceDB server to check connectivity"),
		// All tools have a url argument
		mcp.WithString("url",
			mcp.Description("The URL of the DiceDB server in format 'host:port'"),
			mcp.DefaultString("localhost:7379"),
		),
	)

	echoTool := mcp.NewTool("echo",
		mcp.WithDescription("Echo a message through the DiceDB server"),
		mcp.WithString("url",
			mcp.Description("The URL of the DiceDB server in format 'host:port'"),
			mcp.DefaultString("localhost:7379"),
		),
		mcp.WithString("message",
			mcp.Required(),
			mcp.Description("The message to echo"),
		),
	)

	// Add the pingTool to the server
	s.AddTool(pingTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Get URL with fallback to default
		var url string = "localhost:7379" // Default fallback
		if urlArg, ok := request.Params.Arguments["url"]; ok && urlArg != nil {
			if urlStr, ok := urlArg.(string); ok && urlStr != "" {
				url = urlStr
			}
		}

		host, port := parseHostAndPort(url)

		client, err := dicedb.NewClient(host, port)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("Error connecting to DiceDB: %v", err)), nil
		}

		resp := client.Fire(&wire.Command{Cmd: "PING"})

		return mcp.NewToolResultText(fmt.Sprintf("Response from DiceDB: %s", formatDiceDBResponse(resp))), nil
	})

	s.AddTool(echoTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		message, ok := request.Params.Arguments["message"].(string)
		if !ok || message == "" {
			return mcp.NewToolResultText("Error: message parameter is required"), nil
		}

		var url string = "localhost:7379"
		if urlArg, ok := request.Params.Arguments["url"]; ok && urlArg != nil {
			if urlStr, ok := urlArg.(string); ok && urlStr != "" {
				url = urlStr
			}
		}

		host, port := parseHostAndPort(url)

		client, err := dicedb.NewClient(host, port)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("Error connecting to DiceDB: %v", err)), nil
		}

		resp := client.Fire(&wire.Command{
			Cmd:  "ECHO",
			Args: []string{message},
		})

		return mcp.NewToolResultText(fmt.Sprintf("DiceDB echoed: %s", formatDiceDBResponse(resp))), nil
	})

	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}

// parseHostAndPort splits a URL string in format "host:port" and returns the host and port
func parseHostAndPort(url string) (string, int) {
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

// formatDiceDBResponse formats the DiceDB response in a human-readable way
func formatDiceDBResponse(resp *wire.Response) string {
	if resp.Err != "" {
		return fmt.Sprintf("Error: %s", resp.Err)
	}

	// Handle different value types
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
