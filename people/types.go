package main


// name,surname,gender,birthDate,birthPlace,birthCountry,job
type Person struct {
    Name string `csv:"name"`
    Surname     string `csv:"surname"`
    Gender      string    `csv:"gender"`
	BirthDate string `csv:"birthDate"`
	BirthPlace string `csv:"birthPlace"`
	BirthCountry string `csv:"birthCountry"`
	Job string `csv:"job"`
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

type AnswerType struct {
	Task   string                 `json:"task"`
	ApiKey string                 `json:"apikey"`
	Answer []PersonJSON `json:"answer"`
}