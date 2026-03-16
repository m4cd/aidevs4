package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/m4cd/aidevs4/internal/answer"
	"github.com/m4cd/aidevs4/internal/dates"
	"github.com/m4cd/aidevs4/internal/files"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"

	"github.com/gocarina/gocsv"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error while loading .env file.")
		return
	}

	// Variables
	ApiKey := os.Getenv("API_KEY")
	Url_centrala := os.Getenv("URL_CNTRL")
	Url_verify := Url_centrala + "verify"
	MinAge := 20
	MaxAge := 40

	DataPath := "data/people"
	DataFilename := "people.csv"
	Data := DataPath + "/" + DataFilename
	DataFiltered := DataPath + "/" + "filteredPeople.csv"

	OpenAiApiKey := os.Getenv("OPENAI_API_KEY")
	openAiClient := openai.NewClient(
		option.WithAPIKey(OpenAiApiKey),
	)

	// Downloading people.csv
	files.DownloadFile(DataPath, DataFilename, Url_centrala+"data/"+ApiKey+"/people.csv")

	// Local processing
	csvFile, csvFileError := os.OpenFile(Data, os.O_RDWR, os.ModePerm)

	if csvFileError != nil {
		fmt.Println("CSV file error.")
		return
	}

	defer csvFile.Close()

	var people []*Person

	if unmarshalError := gocsv.UnmarshalFile(csvFile, &people); unmarshalError != nil {
		fmt.Println("CSV unmarshalling error.")
		return
	}

	var filteredPeople []*Person

	for _, person := range people {
		personsAge, err := dates.CalculateAgeColonYYYYMMDD(person.BirthDate)
		if err != nil {
			fmt.Println("Error while calculating age.")
			return
		}
		if person.Gender == "M" && person.BirthPlace == "Grudziądz" && personsAge >= MinAge && personsAge <= MaxAge {
			filteredPeople = append(filteredPeople, person)
		}
	}

	file, err := os.Create(DataFiltered)
	if err != nil {
		fmt.Println("Filtered people CSV error.")
		return
	}
	defer file.Close()
	gocsv.MarshalFile(&filteredPeople, file)

	// OPENAI
	systemMessage := `Z podanego pliku CSV musisz odpowiednio otagować. Mamy do dyspozycji następujące tagi:
- IT
- transport
- edukacja
- medycyna
- praca z ludźmi
- praca z pojazdami
- praca fizyczna
Jedna osoba może mieć wiele tagów.
Odpowiedź zwróć w formacie JSON w polu "answer":
{
       "answer": [
         {
           "name": "Jan",
           "surname": "Kowalski",
           "gender": "M",
           "born": 1987,
           "city": "Warszawa",
           "tags": ["tag1", "tag2"]
         },
         {
           "name": "Anna",
           "surname": "Nowak",
           "gender": "F",
           "born": 1993,
           "city": "Grudziądz",
           "tags": ["tagA", "tagB", "tagC"]
        }
    ]
}
	`

	userMessage := files.ReadFileToString(DataFiltered)

	schema := GenerateSchema()
	completion, err := openAiClient.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
		Model: openai.ChatModelGPT4oMini,
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(systemMessage),
			openai.UserMessage(userMessage),
		},
		ResponseFormat: openai.ChatCompletionNewParamsResponseFormatUnion{
			OfJSONSchema: &schema,
		},
	})

	if err != nil {
		fmt.Println("Chat completion error.")
		return
	}

	// Check for refusal
	choice := completion.Choices[0]
	if choice.Message.Refusal != "" {
		fmt.Println("Refusal error.")
		return
	}

	var result PeopleListJSON
	if err := json.Unmarshal([]byte(choice.Message.Content), &result); err != nil {
		fmt.Println("Unmarshalling error.")
		return
	}

	ans := AnswerType{
		Task:   "people",
		ApiKey: ApiKey,
	}

	for _, r := range result.People {
		for _, i := range r.Tags {
			if i == "transport" {
				ans.Answer = append(ans.Answer, r)
				break
			}
		}

	}

	response := answer.SendAnswer(ans, Url_verify)
	fmt.Println(string(response))

}

func GenerateSchema() openai.ResponseFormatJSONSchemaParam {
	return openai.ResponseFormatJSONSchemaParam{
		JSONSchema: openai.ResponseFormatJSONSchemaJSONSchemaParam{
			Name:   "people_list",
			Strict: openai.Bool(true),
			Schema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"people": map[string]any{
						"type": "array",
						"items": map[string]any{
							"type": "object",
							"properties": map[string]any{
								"name":    map[string]any{"type": "string"},
								"surname": map[string]any{"type": "string"},
								"gender":  map[string]any{"type": "string", "enum": []string{"M", "F"}},
								"born":    map[string]any{"type": "integer"},
								"city":    map[string]any{"type": "string"},
								"tags": map[string]any{
									"type":  "array",
									"items": map[string]any{"type": "string"},
								},
							},
							"required":             []string{"name", "surname", "gender", "born", "city", "tags"},
							"additionalProperties": false,
						},
					},
				},
				"required":             []string{"people"},
				"additionalProperties": false,
			},
		},
	}
}
