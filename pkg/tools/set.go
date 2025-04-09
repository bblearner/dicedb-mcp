package tools

import (
	"context"
	"fmt"
	"strconv"

	"github.com/dicedb/dicedb-go"
	"github.com/dicedb/dicedb-go/wire"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/pottekkat/dicedb-mcp/pkg/utils"
)

// NewSetTool creates a new SET tool for DiceDB
func NewSetTool() mcp.Tool {
	return mcp.NewTool("set",
		mcp.WithDescription("Set a key-value pair in DiceDB"),
		mcp.WithString("url",
			mcp.Description("The URL of the DiceDB server in format 'host:port'"),
			mcp.DefaultString("localhost:7379"),
		),
		mcp.WithString("key",
			mcp.Required(),
			mcp.Description("The key to set"),
		),
		mcp.WithString("value",
			mcp.Required(),
			mcp.Description("The value to set"),
		),
		mcp.WithNumber("ex",
			mcp.Description("Set the expiration time in seconds"),
		),
		mcp.WithNumber("px",
			mcp.Description("Set the expiration time in milliseconds"),
		),
		mcp.WithNumber("exat",
			mcp.Description("Set the expiration time in seconds since epoch"),
		),
		mcp.WithNumber("pxat",
			mcp.Description("Set the expiration time in milliseconds since epoch"),
		),
		mcp.WithBoolean("xx",
			mcp.Description("Only set the key if it already exists"),
		),
		mcp.WithBoolean("nx",
			mcp.Description("Only set the key if it does not already exist"),
		),
		mcp.WithBoolean("keepttl",
			mcp.Description("Keep the existing TTL of the key"),
		),
		mcp.WithBoolean("get",
			mcp.Description("Return the value of the key after setting it"),
		),
	)
}

// HandleSetTool handles the SET tool request
func HandleSetTool(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	key, ok := request.Params.Arguments["key"].(string)
	if !ok || key == "" {
		return nil, fmt.Errorf("missing or empty key parameter")
	}

	value, ok := request.Params.Arguments["value"].(string)
	if !ok {
		return nil, fmt.Errorf("missing value parameter")
	}

	var url string = "localhost:7379"
	if urlArg, ok := request.Params.Arguments["url"]; ok && urlArg != nil {
		if urlStr, ok := urlArg.(string); ok && urlStr != "" {
			url = urlStr
		}
	}

	args := []string{key, value}

	if ex, ok := request.Params.Arguments["ex"].(float64); ok {
		args = append(args, "EX", strconv.FormatInt(int64(ex), 10))
	}

	if px, ok := request.Params.Arguments["px"].(float64); ok {
		args = append(args, "PX", strconv.FormatInt(int64(px), 10))
	}

	if exat, ok := request.Params.Arguments["exat"].(float64); ok {
		args = append(args, "EXAT", strconv.FormatInt(int64(exat), 10))
	}

	if pxat, ok := request.Params.Arguments["pxat"].(float64); ok {
		args = append(args, "PXAT", strconv.FormatInt(int64(pxat), 10))
	}

	if xx, ok := request.Params.Arguments["xx"].(bool); ok && xx {
		args = append(args, "XX")
	}

	if nx, ok := request.Params.Arguments["nx"].(bool); ok && nx {
		args = append(args, "NX")
	}

	if keepttl, ok := request.Params.Arguments["keepttl"].(bool); ok && keepttl {
		args = append(args, "KEEPTTL")
	}

	if get, ok := request.Params.Arguments["get"].(bool); ok && get {
		args = append(args, "GET")
	}

	host, port := utils.ParseHostAndPort(url)
	client, err := dicedb.NewClient(host, port)
	if err != nil {
		return nil, fmt.Errorf("error connecting to DiceDB: %w", err)
	}

	resp := client.Fire(&wire.Command{
		Cmd:  "SET",
		Args: args,
	})

	if resp.Err != "" {
		return nil, fmt.Errorf("DiceDB error: %s", resp.Err)
	}

	formattedResponse := utils.FormatDiceDBResponse(resp)

	var resultMessage string
	if formattedResponse == "OK" {
		resultMessage = fmt.Sprintf("Successfully set key '%s'", key)
	} else if formattedResponse == "(nil)" {
		resultMessage = fmt.Sprintf("Key '%s' was not set (condition not met)", key)
	} else {
		// Display the value of the key also
		resultMessage = fmt.Sprintf("Set key '%s' with value: %s", key, formattedResponse)
	}

	return mcp.NewToolResultText(resultMessage), nil
}
