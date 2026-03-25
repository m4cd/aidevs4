package main

import (
	"github.com/m4cd/aidevs4/internal/coordinates"
	"github.com/openai/openai-go"
)

var tools = []openai.ChatCompletionToolParam{
	{
		Type: "function",
		Function: openai.FunctionDefinitionParam{
			Name:        "get_powerplant_coordinates",
			Description: openai.String("Returns latitude and longitude for a given powerplant identified by the name of a city where it's located."),
			Parameters: openai.FunctionParameters{
				"type": "object",
				"properties": map[string]interface{}{
					"city": map[string]string{
						"type":        "string",
						"description": "The name of the city",
					},
				},
				"required": []string{"city"},
			},
		},
	},
	{
		Type: "function",
		Function: openai.FunctionDefinitionParam{
			Name:        "nearest_powerplant",
			Description: openai.String("Returns the nearest powerplant to any of the coordinates the suspect was seen at and the distance."),
			Parameters: openai.FunctionParameters{
				"type": "object",
				"properties": map[string]interface{}{
					"coordinates": map[string]interface{}{
						"type":        "array",
						"description": "List of coordinates where the suspect was seen",
						"items": map[string]interface{}{
							"type": "object",
							"properties": map[string]interface{}{
								"latitude":  map[string]string{"type": "number", "description": "Latitude"},
								"longitude": map[string]string{"type": "number", "description": "Longitude"},
							},
							"required": []string{"latitude", "longitude"},
						},
					},
					"powerplant_coordinates": map[string]interface{}{
						"type":        "array",
						"description": "List of coordinates of powerplants",
						"items": map[string]interface{}{
							"type": "object",
							"properties": map[string]interface{}{
								"city":      map[string]string{"type": "string", "description": "City"},
								"latitude":  map[string]string{"type": "number", "description": "Latitude"},
								"longitude": map[string]string{"type": "number", "description": "Longitude"},
							},
							"required": []string{"latitude", "longitude"},
						},
					},
				},
				"required": []string{"coordinates", "powerplant_coordinates"},
			},
		},
	},
	{
		Type: "function",
		Function: openai.FunctionDefinitionParam{
			Name:        "get_suspects_locations",
			Description: openai.String("Retrieves locations suspect was seen at."),
			Parameters: openai.FunctionParameters{
				"type": "object",
				"properties": map[string]interface{}{
					"Name":    map[string]string{"type": "string", "description": "Suspects name"},
					"Surname": map[string]string{"type": "string", "description": "Suspects surname"},
				},
				"required": []string{"Name", "Surname"},
			},
		},
	},
	{
		Type: "function",
		Function: openai.FunctionDefinitionParam{
			Name:        "send_answer",
			Description: openai.String("Sends answer."),
			Parameters: openai.FunctionParameters{
				"type": "object",
				"properties": map[string]interface{}{
					"Name":        map[string]string{"type": "string", "description": "Suspects name"},
					"Surname":     map[string]string{"type": "string", "description": "Suspects surname"},
					"AccessLevel": map[string]string{"type": "integer", "description": "Suspects accesslevel"},
					"PowerPlant":  map[string]string{"type": "string", "description": "Code of the Powerplant"},
				},
				"required": []string{"Name", "Surname", "AccessLevel", "PowerPlant"},
			},
		},
	},
	{
		Type: "function",
		Function: openai.FunctionDefinitionParam{
			Name:        "get_accesslevel",
			Description: openai.String("Return suspects name, surname and access level"),
			Parameters: openai.FunctionParameters{
				"type": "object",
				"properties": map[string]interface{}{
					"Name":      map[string]string{"type": "string", "description": "Suspects name"},
					"Surname":   map[string]string{"type": "string", "description": "Suspects surname"},
					"Birthdate": map[string]string{"type": "string", "description": "Suspects date of birth"},
				},
				"required": []string{"Name", "Surname", "Birthdate"},
			},
		},
	},
	{
		Type: "function",
		Function: openai.FunctionDefinitionParam{
			Name:        "success",
			Description: openai.String("Success function ending agent loop"),
			Parameters: openai.FunctionParameters{
				"type":       "object",
				"properties": map[string]interface{}{},
				"required":   []string{},
			},
		},
	},
}

func resolveCoordinates(city string) *coordinates.Coordinate {
	coords := map[string]coordinates.Coordinate{
		"Zabrze":               {Latitude: 50.3249, Longitude: 18.7857},
		"Piotrków Trybunalski": {Latitude: 51.4058, Longitude: 19.7031},
		"Grudziądz":            {Latitude: 53.4836, Longitude: 18.7536},
		"Tczew":                {Latitude: 54.0922, Longitude: 18.7787},
		"Radom":                {Latitude: 51.4027, Longitude: 21.1471},
		"Chelmno":              {Latitude: 53.3494, Longitude: 18.4247},
		"Żarnowiec":            {Latitude: 54.6667, Longitude: 18.1667},
	}
	if c, ok := coords[city]; ok {
		return &c
	}
	return &coordinates.Coordinate{
		Latitude:  0,
		Longitude: 0,
	}
}

type PowerplantCoordinate struct {
	City      string  `json:"city"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type nearestPowerplantInput struct {
	Coordinates           []coordinates.Coordinate `json:"coordinates"`
	PowerplantCoordinates []PowerplantCoordinate   `json:"powerplant_coordinates"`
}

type nearestPowerplantStruct struct {
	distance float64
	city     string
}

func nearestPowerplant(input nearestPowerplantInput) nearestPowerplantStruct {

	result := nearestPowerplantStruct{
		distance: 1000,
		city:     "",
	}

	for _, powerplant := range input.PowerplantCoordinates {
		for _, suspectLocation := range input.Coordinates {
			distance := coordinates.Haversine(suspectLocation.Latitude, suspectLocation.Longitude, powerplant.Latitude, powerplant.Longitude)
			if distance < result.distance {
				result.distance = distance
				result.city = powerplant.City
			}
		}
	}
	return result
}
