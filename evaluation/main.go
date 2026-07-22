package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
	"github.com/m4cd/aidevs4/internal/answer"
	"github.com/m4cd/aidevs4/internal/files"
	"github.com/m4cd/aidevs4/internal/types"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
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
	Task := "evaluation"

	OpenAiApiKey := os.Getenv("OPENAI_API_KEY")
	openAiClient := openai.NewClient(
		option.WithAPIKey(OpenAiApiKey),
	)

	jsonFiles, err := filepath.Glob("data/" + Task + "/*.json")
	if err != nil {
		log.Fatal(err)
	}

	wrongData := []string{}
	operatorNotes := []string{}

	for _, path := range jsonFiles {

		name := filepath.Base(path)

		raw, err := os.ReadFile(path)
		if err != nil {
			log.Printf("read %s: %v", path, err)
			continue
		}

		var reading types.EvaluationSensorReading
		if err := json.Unmarshal(raw, &reading); err != nil {
			log.Printf("parse %s: %v", path, err)
			continue
		}
		// fmt.Println(reading.String())
		if problems := types.Validate(reading); len(problems) > 0 {
			log.Printf("%s: %s", path, strings.Join(problems, "; "))
			wrongData = append(wrongData, name)
		} else {
			firstPart := strings.Split(reading.OperatorNotes, ",")[0]

			nonexistent := true
			for _, el := range operatorNotes {
				if el == firstPart {
					nonexistent = false
					break
				}
			}

			if nonexistent {
				operatorNotes = append(operatorNotes, firstPart)
			}

		}

	}

	messages := []openai.ChatCompletionMessageParamUnion{}
	systemPrompt := files.ReadFileToString(Task + "/" + "prompt.md")

	userMessage := ""
	for _, on := range operatorNotes {
		userMessage = userMessage + on + "\n"
	}

	messages = append(messages, openai.SystemMessage(systemPrompt))
	messages = append(messages, openai.UserMessage(userMessage))

	responseOpenAI, err := openAiClient.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
		Model:    "gpt-5.4-mini",
		Messages: messages,
	})

	if err != nil {
		fmt.Println("Chat completion error.")
		fmt.Println(err)
		return
	}

	choice := responseOpenAI.Choices[0]
	llmtxt := choice.Message.Content

	llmSelectedtxts := strings.Split(llmtxt, "\n")

	for _, path := range jsonFiles {

		name := filepath.Base(path)

		raw, err := os.ReadFile(path)
		if err != nil {
			log.Printf("read %s: %v", path, err)
			continue
		}

		var reading types.EvaluationSensorReading
		if err := json.Unmarshal(raw, &reading); err != nil {
			log.Printf("parse %s: %v", path, err)
			continue
		}
		firstPart := strings.Split(reading.OperatorNotes, ",")[0]
		for _, lstxts := range llmSelectedtxts {
			if firstPart == lstxts {

				wrongData = append(wrongData, name)

			}
		}
	}

	var ans types.AnswerS03E01
	ans.ApiKey = ApiKey
	ans.Task = Task
	ans.Answer.Recheck = wrongData

	response := answer.SendAnswerReturnHttp(ans, Url_verify)

	answer.PrintResponse(response)
}
