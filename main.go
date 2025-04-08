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
		mcp.WithString("url",
			mcp.Description("The URL of the DiceDB server in format 'host:port'"),
			mcp.DefaultString("localhost:7379"),
		),
	)

	s.AddTool(pingTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Extract the URL from the request arguments
		url := request.Params.Arguments["url"].(string)

		// Parse host and port from URL
		host, port := parseHostAndPort(url)

		// Create a new DiceDB client
		client, err := dicedb.NewClient(host, port)
		if err != nil {
			return mcp.NewToolResultText(fmt.Sprintf("Error connecting to DiceDB: %v", err)), nil
		}

		// Send PING command
		resp := client.Fire(&wire.Command{Cmd: "PING"})

		// Return the response to the MCP client
		return mcp.NewToolResultText(fmt.Sprintf("Response from DiceDB: %v", resp)), nil
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
