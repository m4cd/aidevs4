package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/joho/godotenv"
	"github.com/m4cd/aidevs4/internal/answer"
	"github.com/m4cd/aidevs4/internal/files"
	"github.com/m4cd/aidevs4/internal/llm"
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
	Task := "categorize"
	csvFilename := "categorize.csv"
	csvURL := Url_centrala + "data/" + ApiKey + "/" + csvFilename
	csvPath := "data/categorize"

	Answer := types.AnswerS02E01{
		ApiKey: ApiKey,
		Task:   Task,
	}

	OpenAiApiKey := os.Getenv("OPENAI_API_KEY")
	openAiClient := openai.NewClient(
		option.WithAPIKey(OpenAiApiKey),
	)

	messages := []openai.ChatCompletionMessageParamUnion{}
	// var response http.Response

	systemPrompt := files.ReadFileToString(Task + "/" + "prompt.md")
	messages = append(messages, openai.SystemMessage(systemPrompt))

	for {

		// fmt.Println("[+] Printing all messages at the beginning of the loop...")
		// llm.PrintAllMessages(messages)

		responseOpenAI, err := openAiClient.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
			Model:    "gpt-5.4",
			Tools:    tools,
			Messages: messages,
		})
		if err != nil {
			fmt.Println("Chat completion error.")
			fmt.Println(err)
			return
		}

		choice := responseOpenAI.Choices[0]

		if choice.Message.Content != "" {
			fmt.Println(choice.Message.Content)
		} else {

			messages = append(messages, openai.ChatCompletionMessageParamUnion{
				OfAssistant: &openai.ChatCompletionAssistantMessageParam{
					ToolCalls: llm.ToToolCallParams(choice.Message.ToolCalls),
				},
			})

			for _, toolCall := range choice.Message.ToolCalls {
				fmt.Println("[+] Function " + toolCall.Function.Name + " chosen...")

				switch toolCall.Function.Name {
				case "reset":
					var resetPrompt types.AnswerActionS02E01
					resetPrompt.Prompt = "reset"

					Answer.Answer = resetPrompt

					response := answer.SendAnswerReturnHttp(Answer, Url_verify)

					body, err := io.ReadAll(response.Body)
					if err != nil {
						fmt.Printf("Error reading body: %v\n", err)
						return
					}

					b, _ := json.Marshal(body)
					var rawBase64 string
					json.Unmarshal(b, &rawBase64)
					resultBytes, _ := base64.StdEncoding.DecodeString(rawBase64)
					result := string(resultBytes)
					fmt.Printf("Result: %s\n", result)

					messages = append(messages, openai.ToolMessage(result, toolCall.ID))

				case "download":
					// fmt.Println("reset")
					fmt.Println("[+] Downloading CSV...")
					files.DownloadFile(csvPath, csvFilename, csvURL)
					fmt.Println(csvPath)
					csvContents := files.ReadFileToString(csvPath + "/" + csvFilename)
					fmt.Println(csvContents)

					messages = append(messages, openai.ToolMessage(csvContents, toolCall.ID))

				case "prompt":
					var input types.AnswerActionS02E01

					json.Unmarshal([]byte(toolCall.Function.Arguments), &input)

					Answer.Answer = input

					response := answer.SendAnswerReturnHttp(Answer, Url_verify)

					body, err := io.ReadAll(response.Body)
					if err != nil {
						fmt.Printf("Error reading body: %v\n", err)
						return
					}

					b, _ := json.Marshal(body)
					var rawBase64 string
					json.Unmarshal(b, &rawBase64)
					resultBytes, _ := base64.StdEncoding.DecodeString(rawBase64)
					result := string(resultBytes)
					fmt.Printf("Result: %s\n", result)

					messages = append(messages, openai.ToolMessage(result, toolCall.ID))
				case "success":
					var input types.SuccessInput
					json.Unmarshal([]byte(toolCall.Function.Arguments), &input)

					fmt.Println(input.Flag)

					return
				}

				fmt.Println()
			}

		}

	}
}
