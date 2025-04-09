package tools

import (
	"context"
	"fmt"

	"github.com/dicedb/dicedb-go"
	"github.com/dicedb/dicedb-go/wire"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/pottekkat/dicedb-mcp/pkg/utils"
)

// NewEchoTool creates a new echo tool for DiceDB
func NewEchoTool() mcp.Tool {
	return mcp.NewTool("echo",
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
}

// HandleEchoTool handles the echo tool request
func HandleEchoTool(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	message, ok := request.Params.Arguments["message"].(string)
	if !ok || message == "" {
		return nil, fmt.Errorf("missing or empty message parameter")
	}

	var url string = "localhost:7379"
	if urlArg, ok := request.Params.Arguments["url"]; ok && urlArg != nil {
		if urlStr, ok := urlArg.(string); ok && urlStr != "" {
			url = urlStr
		}
	}

	host, port := utils.ParseHostAndPort(url)

	client, err := dicedb.NewClient(host, port)
	if err != nil {
		return nil, fmt.Errorf("error connecting to DiceDB: %w", err)
	}

	resp := client.Fire(&wire.Command{
		Cmd:  "ECHO",
		Args: []string{message},
	})

	// Check if DiceDB returned an error
	if resp.Err != "" {
		return nil, fmt.Errorf("DiceDB error: %s", resp.Err)
	}

	return mcp.NewToolResultText(fmt.Sprintf("DiceDB echoed: %s", utils.FormatDiceDBResponse(resp))), nil
}
