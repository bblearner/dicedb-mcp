package main

import (
	"fmt"

	"github.com/mark3labs/mcp-go/server"
	"github.com/pottekkat/dicedb-mcp/pkg/tools"
)

func main() {
	// Create a new MCP server
	s := server.NewMCPServer(
		"DiceDB MCP",
		"0.1.0",
		server.WithResourceCapabilities(true, true),
		server.WithLogging(),
	)

	// Create and add the ping tool
	pingTool := tools.NewPingTool()
	s.AddTool(pingTool, tools.HandlePingTool)

	echoTool := tools.NewEchoTool()
	s.AddTool(echoTool, tools.HandleEchoTool)

	// Start the server
	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}
