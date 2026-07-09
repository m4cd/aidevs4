package main

import "github.com/openai/openai-go"

var tools = []openai.ChatCompletionToolParam{
	{
		Type: "function",
		Function: openai.FunctionDefinitionParam{
			Name:        "help",
			Description: openai.String("Show available actions and parameters of the API."),
			Parameters: openai.FunctionParameters{
				"type":       "object",
				"properties": map[string]interface{}{},
				"required":   []string{},
			},
		},
	},
	{
		Type: "function",
		Function: openai.FunctionDefinitionParam{
			Name:        "getInbox",
			Description: openai.String("Return list of threads in your mailbox."),
			Parameters: openai.FunctionParameters{
				"type": "object",
				"properties": map[string]interface{}{
					"page": map[string]interface{}{
						"type":        "integer",
						"minimum":     1,
						"default":     1,
						"description": "Optional. Page number, integer >= 1. Default: 1.",
					},
					"perPage": map[string]interface{}{
						"type":        "integer",
						"minimum":     5,
						"maximum":     20,
						"default":     5,
						"description": "Optional. Number of threads per page, integer between 5 and 20. Default: 5.",
					},
				},
				"required": []string{"page", "perPage"},
			},
		},
	},
	{
		Type: "function",
		Function: openai.FunctionDefinitionParam{
			Name:        "getThread",
			Description: openai.String("Return rowID and messageID list for a selected thread. No message body."),
			Parameters: openai.FunctionParameters{
				"type": "object",
				"properties": map[string]interface{}{
					"threadID": map[string]interface{}{
						"type":        "integer",
						"description": "Numeric thread identifier.",
					},
				},
				"required": []string{"threadID"},
			},
		},
	},
	{
		Type: "function",
		Function: openai.FunctionDefinitionParam{
			Name:        "getMessages",
			Description: openai.String("Return one or more messages by rowID/messageID (hash)."),
			Parameters: openai.FunctionParameters{
				"type": "object",
				"properties": map[string]interface{}{
					"ids": map[string]interface{}{
						"description": "Numeric rowID, 32-char messageID, or an array of them.",
						"anyOf": []map[string]interface{}{
							{"type": "integer"},
							{"type": "string"},
							{
								"type":  "array",
								"items": map[string]interface{}{"type": []string{"integer", "string"}},
							},
						},
					},
				},
				"required": []string{"ids"},
			},
		},
	},
	{
		Type: "function",
		Function: openai.FunctionDefinitionParam{
			Name:        "search",
			Description: openai.String("Search messages with full-text style query and Gmail-like operators."),
			Parameters: openai.FunctionParameters{
				"type": "object",
				"properties": map[string]interface{}{
					"query": map[string]interface{}{
						"type":        "string",
						"description": `Supports words, "phrase", -exclude, from:, to:, subject:, subject:"phrase", subject:(phrase), OR, AND. Missing operator means AND.`,
					},
					"page": map[string]interface{}{
						"type":        "integer",
						"minimum":     1,
						"default":     1,
						"description": "Optional. Page number, integer >= 1. Default: 1.",
					},
					"perPage": map[string]interface{}{
						"type":        "integer",
						"minimum":     5,
						"maximum":     20,
						"default":     5,
						"description": "Optional. Number of threads per page, integer between 5 and 20. Default: 5.",
					},
				},
				"required": []string{"query"},
			},
		},
	},
	{
		Type: "function",
		Function: openai.FunctionDefinitionParam{
			Name:        "reset",
			Description: openai.String("Reset request counter for this apikey in memcache (in case of fuckup)."),
			Parameters: openai.FunctionParameters{
				"type":       "object",
				"properties": map[string]interface{}{},
				"required":   []string{},
			},
		},
	},
	{
		Type: "function",
		Function: openai.FunctionDefinitionParam{
			Name:        "verify",
			Description: openai.String("To be called when the flag is found. Report the recovered password, date, and confirmation code."),
			Parameters: openai.FunctionParameters{
				"type": "object",
				"properties": map[string]interface{}{
					"password": map[string]interface{}{
						"type":        "string",
						"description": "The recovered password.",
					},
					"date": map[string]interface{}{
						"type":        "string",
						"description": "The recovered date.",
					},
					"confirmation_code": map[string]interface{}{
						"type":        "string",
						"description": "The recovered confirmation code.",
					},
				},
				"required": []string{"password", "date", "confirmation_code"},
			},
		},
	},
}
