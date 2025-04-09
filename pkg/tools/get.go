package tools

import (
	"context"
	"fmt"

	"github.com/dicedb/dicedb-go"
	"github.com/dicedb/dicedb-go/wire"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/pottekkat/dicedb-mcp/pkg/utils"
)

// NewGetTool creates a new GET tool for DiceDB
func NewGetTool() mcp.Tool {
	return mcp.NewTool("get",
		mcp.WithDescription("Get a value from DiceDB by key"),
		mcp.WithString("url",
			mcp.Description("The URL of the DiceDB server in format 'host:port'"),
			mcp.DefaultString("localhost:7379"),
		),
		mcp.WithString("key",
			mcp.Required(),
			mcp.Description("The key to retrieve from DiceDB"),
		),
	)
}

// HandleGetTool handles the GET tool request
func HandleGetTool(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	key, ok := request.Params.Arguments["key"].(string)
	if !ok || key == "" {
		return mcp.NewToolResultText("Error: key parameter is required"), nil
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
		return mcp.NewToolResultText(fmt.Sprintf("Error connecting to DiceDB: %v", err)), nil
	}

	resp := client.Fire(&wire.Command{
		Cmd:  "GET",
		Args: []string{key},
	})

	value := utils.FormatDiceDBResponse(resp)

	// If value is "(nil)", the key doesn't exist
	if value == "(nil)" {
		return mcp.NewToolResultText(fmt.Sprintf("Key '%s' not found", key)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Value for key '%s': %s", key, value)), nil
}
