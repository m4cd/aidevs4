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
			Name:        "reconfigure",
			Description: openai.String("Enable reconfigure mode for the given route. Must be called before setstatus."),
			Parameters: openai.FunctionParameters{
				"type": "object",
				"properties": map[string]interface{}{
					"route": map[string]string{
						"type":        "string",
						"description": "Route identifier in format [a-z]-[0-9]{1,2}, e.g. 'a-1' or 'b-12' (case-insensitive)",
					},
				},
				"required": []string{"route"},
			},
		},
	},
	{
		Type: "function",
		Function: openai.FunctionDefinitionParam{
			Name:        "getstatus",
			Description: openai.String("Get current status for the given route."),
			Parameters: openai.FunctionParameters{
				"type": "object",
				"properties": map[string]interface{}{
					"route": map[string]string{
						"type":        "string",
						"description": "Route identifier in format [a-z]-[0-9]{1,2}, e.g. 'a-1' or 'b-12' (case-insensitive)",
					},
				},
				"required": []string{"route"},
			},
		},
	},
	{
		Type: "function",
		Function: openai.FunctionDefinitionParam{
			Name:        "setstatus",
			Description: openai.String("Set route status while in reconfigure mode. Allowed values: RTOPEN (open the route), RTCLOSE (close the route). Requires reconfigure to be called first."),
			Parameters: openai.FunctionParameters{
				"type": "object",
				"properties": map[string]interface{}{
					"route": map[string]string{
						"type":        "string",
						"description": "Route identifier in format [a-z]-[0-9]{1,2}, e.g. 'a-1' or 'b-12' (case-insensitive)",
					},
					"value": map[string]interface{}{
						"type":        "string",
						"description": "New status for the route",
						"enum":        []string{"RTOPEN", "RTCLOSE"},
					},
				},
				"required": []string{"route", "value"},
			},
		},
	},
	{
		Type: "function",
		Function: openai.FunctionDefinitionParam{
			Name:        "save",
			Description: openai.String("Exit reconfigure mode for the given route, saving any status changes made."),
			Parameters: openai.FunctionParameters{
				"type": "object",
				"properties": map[string]interface{}{
					"route": map[string]string{
						"type":        "string",
						"description": "Route identifier in format [a-z]-[0-9]{1,2}, e.g. 'a-1' or 'b-12' (case-insensitive)",
					},
				},
				"required": []string{"route"},
			},
		},
	},
	{
		Type: "function",
		Function: openai.FunctionDefinitionParam{
			Name:        "success",
			Description: openai.String("To be called when flag is found."),
			Parameters: openai.FunctionParameters{
				"type": "object",
				"properties": map[string]interface{}{
					"flag": map[string]string{
						"type":        "string",
						"description": "The found result",
					},
				},
				"required": []string{"flag"},
			},
		},
	},
}
