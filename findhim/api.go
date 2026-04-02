package main

import (
	"fmt"

	"github.com/m4cd/aidevs4/internal/answer"
	"github.com/m4cd/aidevs4/internal/coordinates"
	"github.com/m4cd/aidevs4/internal/dates"
	"github.com/m4cd/aidevs4/internal/strings"
	"github.com/m4cd/aidevs4/internal/types"
)

type LocationApiCallInput struct {
	Name    string `json:"name"`
	Surname string `json:"surname"`
}

func LocationApiCall(key string, url string, name string, surname string) []coordinates.Coordinate {
	var suspectPost types.LocationPostJson
	suspectPost.Name = name
	suspectPost.Surname = surname
	suspectPost.ApiKey = key

	CandidatesLocations, err := answer.SendPostJson(url, &suspectPost)
	if err != nil {
		fmt.Println("Error sending candidate to location api.")
		return nil
	}
	CoordinatesArray, err := coordinates.UnmarshalCoordinates(CandidatesLocations)
	if err != nil {
		fmt.Println("Error running 'coordinates.UnmarshalCoordinates(CandidatesLocations)'.")
		return nil
	}
	return CoordinatesArray
}

type AccessLevelApiCallInput struct {
	Name      string `json:"name"`
	Surname   string `json:"surname"`
	Birthdate string `json:"birthdate"`
}

func AccessLevelApiCall(key string, url string, name string, surname string, birthdate string) types.AccessLevelResponse {
	var Accesslevel types.AccesslevelPostJson
	Accesslevel.Name = name
	Accesslevel.Surname = surname
	Accesslevel.ApiKey = key

	var err error
	Accesslevel.BirthYear, err = dates.ExtractYearYYYYMMDD(birthdate)
	if err != nil {
		fmt.Println("Error extracting year from birthdate string.")
		return types.AccessLevelResponse{}
	}

	AccesslevelResponseString, err := answer.SendPostJson(url, &Accesslevel)
	if err != nil {
		fmt.Println("Error sending candidate to accesslevel api.")
		return types.AccessLevelResponse{}
	}
	AccesslevelResponse, _ := strings.UnmarshalJSON[types.AccessLevelResponse](AccesslevelResponseString)
	AccesslevelResponse.Print()
	return *AccesslevelResponse
}
