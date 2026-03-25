package types

import (
	"fmt"
	"reflect"
)

// S01E01
type Person struct {
	Name         string `csv:"name"`
	Surname      string `csv:"surname"`
	Gender       string `csv:"gender"`
	BirthDate    string `csv:"birthDate"`
	BirthPlace   string `csv:"birthPlace"`
	BirthCountry string `csv:"birthCountry"`
	Job          string `csv:"job"`
}

func (p *Person) Print() {
	v := reflect.ValueOf(*p)
	t := reflect.TypeOf(*p)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)
		fmt.Printf("%s: %s\n", field.Tag.Get("csv"), value)
	}
	fmt.Println()
}


type PersonJSON struct {
	Name    string   `json:"name"`
	Surname string   `json:"surname"`
	Gender  string   `json:"gender"`
	Born    int      `json:"born"`
	City    string   `json:"city"`
	Tags    []string `json:"tags"`
}

type PeopleListJSON struct {
	People []PersonJSON `json:"people"`
}

type AnswerS01E01 struct {
	Task   string       `json:"task"`
	ApiKey string       `json:"apikey"`
	Answer []PersonJSON `json:"answer"`
}
