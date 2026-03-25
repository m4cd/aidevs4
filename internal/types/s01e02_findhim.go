package types

import (
	"fmt"
)

type LocationPostJson struct {
	Name    string `json:"name"`
	Surname string `json:"surname"`
	ApiKey  string `json:"apikey"`
}

type AccesslevelPostJson struct {
	Name      string `json:"name"`
	Surname   string `json:"surname"`
	ApiKey    string `json:"apikey"`
	BirthYear int32  `json:"birthYear"`
}

type PowerPlant struct {
	IsActive bool   `json:"is_active"`
	Power    string `json:"power"`
	Code     string `json:"code"`
}

func (p *PowerPlant) Print() {
	fmt.Printf("Is active:  %v\n", p.IsActive)
	fmt.Printf("Power:      %s\n", p.Power)
	fmt.Printf("Code:       %s\n", p.Code)
}

type PowerPlantsJson struct {
	PowerPlants map[string]PowerPlant `json:"power_plants"`
}

func (pp *PowerPlantsJson) Print() {
	for k, p := range pp.PowerPlants {
		fmt.Printf("Name:       %s\n", k)
		p.Print()
		fmt.Println()
	}
}

type AccessLevelResponse struct {
	Name        string `json:"name"`
	Surname     string `json:"surname"`
	AccessLevel int    `json:"accessLevel"`
}

func (alr *AccessLevelResponse) Print() {
	fmt.Printf("Name: %v\n", alr.Name)
	fmt.Printf("Surname: %s\n", alr.Surname)
	fmt.Printf("AccessLevel: %v\n", alr.AccessLevel)
}

type AnswerSuspectS01E02 struct {
	Name        string `json:"name"`
	Surname     string `json:"surname"`
	AccessLevel int `json:"accessLevel"`
	PowerPlant  string `json:"powerPlant"`
}

func (as *AnswerSuspectS01E02) Print() {
	fmt.Printf("Name: %v\n", as.Name)
	fmt.Printf("Surname: %s\n", as.Surname)
	fmt.Printf("AccessLevel: %v\n", as.AccessLevel)
	fmt.Printf("PowerPlant: %v\n", as.PowerPlant)
}

type AnswerS01E02 struct {
	Task   string              `json:"task"`
	ApiKey string              `json:"apikey"`
	Answer AnswerSuspectS01E02 `json:"answer"`
}

type ResolveCoordinatesInput struct {
	City string `json:"city"`
}
