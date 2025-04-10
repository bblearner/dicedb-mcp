package tools

import (
	"context"
	"fmt"
	"strconv"

	"github.com/dicedb/dicedb-go/wire"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/pottekkat/dicedb-mcp/pkg/utils"
)

// NewDelTool creates a new DEL tool for DiceDB
func NewDelTool() mcp.Tool {
	return mcp.NewTool("del",
		mcp.WithDescription("Delete one or more keys from DiceDB"),
		utils.CommonURLParam(),
		mcp.WithArray("keys",
			mcp.Required(),
			mcp.Description("The keys to delete from DiceDB"),
		),
	)
}

// HandleDelTool handles the DEL tool request
func HandleDelTool(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	keys, ok := request.Params.Arguments["keys"].([]interface{})
	if !ok || len(keys) == 0 {
		return nil, fmt.Errorf("missing or empty keys parameter")
	}

	// Convert the keys to strings
	stringKeys := make([]string, len(keys))
	for i, key := range keys {
		stringKeys[i] = fmt.Sprintf("%v", key)
	}

	client, err := utils.GetClientFromRequest(request)
	if err != nil {
		return nil, err
	}

	resp := client.Fire(&wire.Command{
		Cmd:  "DEL",
		Args: stringKeys,
	})

	if resp.Err != "" {
		return nil, fmt.Errorf("DiceDB error: %s", resp.Err)
	}

	value := utils.FormatDiceDBResponse(resp)

	// The response is the number of keys deleted
	count, err := strconv.Atoi(value)
	if err != nil {
		return nil, fmt.Errorf("unexpected response format: %s", value)
	}

	if count == 0 {
		return mcp.NewToolResultText("No keys were deleted"), nil
	} else if count == 1 {
		return mcp.NewToolResultText("Deleted 1 key"), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Deleted %d keys", count)), nil
}
