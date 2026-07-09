package main

import "github.com/openai/openai-go"

var tools = []openai.ChatCompletionToolParam{
	// {
	// 	Type: "function",
	// 	Function: openai.FunctionDefinitionParam{
	// 		Name:        "mapa",
	// 		Description: openai.String("Dostarcza URL pod którym znajduje się mapa elektrowni wraz z naniesioną siatką."),
	// 		Parameters: openai.FunctionParameters{
	// 			"type":       "object",
	// 			"properties": map[string]interface{}{},
	// 			"required":   []string{},
	// 		},
	// 	},
	// },
	// {
	// 	Type: "function",
	// 	Function: openai.FunctionDefinitionParam{
	// 		Name:        "docs",
	// 		Description: openai.String("Dostarcza URL pod którym znajduję się dokumentacji API drona."),
	// 		Parameters: openai.FunctionParameters{
	// 			"type": "object",
	// 			"properties": map[string]interface{}{
	// 				"route": map[string]string{
	// 					"type":        "string",
	// 					"description": "Route identifier in format [a-z]-[0-9]{1,2}, e.g. 'a-1' or 'b-12' (case-insensitive)",
	// 				},
	// 			},
	// 			"required": []string{"route"},
	// 		},
	// 	},
	// },
	{
		Type: "function",
		Function: openai.FunctionDefinitionParam{
			Name:        "execute",
			Description: openai.String("Wysyła zidentyfikowane instrukcje do drona i zwraca odpowiedź."),
			Parameters: openai.FunctionParameters{
				"type": "object",
				"properties": map[string]interface{}{
					"instructions": map[string]interface{}{
						"type":        "array",
						"description": "Lista instrukcji, które mają zostać wykonane przez drona.",
						"items": map[string]string{
							"type": "string",
						},
					},
				},
				"required": []string{"instructions"},
			},
		},
	},
	{
		Type: "function",
		Function: openai.FunctionDefinitionParam{
			Name:        "flg",
			Description: openai.String("Kończy zadanie. Jako dane wejściowe przyjmuję znalezioną w zadaniu flagę {FLG:...}"),
			Parameters: openai.FunctionParameters{
				"type": "object",
				"properties": map[string]interface{}{
					"flg": map[string]string{
						"type":        "string",
						"description": "Znaleziona flaga",
					},
				},
				"required": []string{"flg"},
			},
		},
	},
}
