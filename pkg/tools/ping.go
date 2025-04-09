package tools

import (
	"context"
	"fmt"

	"github.com/dicedb/dicedb-go"
	"github.com/dicedb/dicedb-go/wire"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/pottekkat/dicedb-mcp/pkg/utils"
)

// NewPingTool creates a new ping tool for DiceDB
func NewPingTool() mcp.Tool {
	return mcp.NewTool("ping",
		mcp.WithDescription("Ping the DiceDB server to check connectivity"),
		// All tools have a url argument
		mcp.WithString("url",
			mcp.Description("The URL of the DiceDB server in format 'host:port'"),
			mcp.DefaultString("localhost:7379"),
		),
	)
}

// HandlePingTool handles the ping tool request
func HandlePingTool(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get URL with fallback to default
	var url string = "localhost:7379" // Default fallback
	if urlArg, ok := request.Params.Arguments["url"]; ok && urlArg != nil {
		if urlStr, ok := urlArg.(string); ok && urlStr != "" {
			url = urlStr
		}
	}

	host, port := utils.ParseHostAndPort(url)

	// Create a new DiceDB client
	client, err := dicedb.NewClient(host, port)
	if err != nil {
		return nil, fmt.Errorf("error connecting to DiceDB: %w", err)
	}

	// Run the PING command on the DiceDB server
	resp := client.Fire(&wire.Command{Cmd: "PING"})

	// Check if DiceDB returned an error
	if resp.Err != "" {
		return nil, fmt.Errorf("DiceDB error: %s", resp.Err)
	}

	// Return the response to the MCP client
	return mcp.NewToolResultText(fmt.Sprintf("Response from DiceDB: %s", utils.FormatDiceDBResponse(resp))), nil
}
