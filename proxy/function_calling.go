package main

import (
	"encoding/json"
	"fmt"

	"github.com/m4cd/aidevs4/internal/answer"
	"github.com/m4cd/aidevs4/internal/structs"
	"github.com/m4cd/aidevs4/internal/types"
	"github.com/openai/openai-go"
)

var tools = []openai.ChatCompletionToolParam{
	{
		Type: "function",
		Function: openai.FunctionDefinitionParam{
			Name:        "check_package",
			Description: openai.String("Zwraca status paczki o danym ID."),
			Parameters: openai.FunctionParameters{
				"type": "object",
				"properties": map[string]any{
					"packageid": map[string]any{
						"type":        "string",
						"description": "ID paczki do sprawdzenia",
					},
				},
				"required": []string{"packageid"},
			},
		},
	},
	{
		Type: "function",
		Function: openai.FunctionDefinitionParam{
			Name:        "redirect_package",
			Description: openai.String("Przekierowuje paczkę o danym ID do określonej destynacji"),
			Parameters: openai.FunctionParameters{
				"type": "object",
				"properties": map[string]any{
					"packageid": map[string]any{
						"type":        "string",
						"description": "ID paczki do przekierowania",
					},
					"destination": map[string]any{
						"type":        "string",
						"description": "Destynacja do której paczka ma zostać przekierowana",
					},
					"code": map[string]any{
						"type":        "string",
						"description": "Kod autoryzacyjny",
					},
				},
				"required": []string{"packageid", "destination", "code"},
			},
		},
	},
}

func RedirectPackage(endpoint string, p types.PackageRedirect) types.PackageRedirectResponse {

	status := types.PackageRedirectResponse{}
	PackageStatusString, err := answer.SendPostJson(endpoint, &p)
	if err != nil {
		fmt.Println("Error checking package status.")
		return status
	}
	err = json.Unmarshal([]byte(PackageStatusString), &status)
	if err != nil {
		fmt.Println("Error unmarshalling PackageStatusString json.")
		return status
	}

	return status
}

type PackageCheckApiCallInput struct {
	PackageID string `json:"packageid"`
}

func PackageCheckApiCall(input PackageCheckApiCallInput, key string, url string) types.PackageCheckResponse {

	Package := types.NewPackageCheck()
	Package.ApiKey = key
	Package.PackageID = input.PackageID


	PackageStatusResponse, err := answer.SendPostJson(url, &Package)
	if err != nil {
		fmt.Println("Error sending Package to check_package api.")
		return types.PackageCheckResponse{}
	}

	var Status types.PackageCheckResponse
	err = json.Unmarshal([]byte(PackageStatusResponse), &Status)
	if err != nil {
		fmt.Println("Error unmarshalling PackageStatusResponse json.")
		return types.PackageCheckResponse{}
	}
	structs.PrintStruct(Status)
	return Status
}

type PackageRedirectApiCallInput struct {
	PackageID string `json:"packageid"`
	Destination string `json:"destination"`
	Code string `json:"code"`
}

func PackageRedirectApiCall(input PackageRedirectApiCallInput, key string, url string) types.PackageRedirectResponse {

	Package := types.NewPackageRedirect()
	Package.ApiKey = key
	Package.PackageID = input.PackageID
	Package.Destination = input.Destination
	Package.Code = input.Code


	PackageRedirectResponse, err := answer.SendPostJson(url, &Package)
	if err != nil {
		fmt.Println("Error sending Package to redirect_package api.")
		return types.PackageRedirectResponse{}
	}

	fmt.Println("[+] PackageRedirectResponse")
	fmt.Println(PackageRedirectResponse)

	var Status types.PackageRedirectResponse
	err = json.Unmarshal([]byte(PackageRedirectResponse), &Status)
	if err != nil {
		fmt.Println("Error unmarshalling PackageRedirectResponse json.")
		return types.PackageRedirectResponse{}
	}
	structs.PrintStruct(Status)
	return Status
}
