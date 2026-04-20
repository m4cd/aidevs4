package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

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
	Task := "railway"
	waitingTime := 5
	maxIterations := 15

	OpenAiApiKey := os.Getenv("OPENAI_API_KEY")
	openAiClient := openai.NewClient(
		option.WithAPIKey(OpenAiApiKey),
	)

	Action := types.HelpInput{
		Action: "help",
	}

	ans := types.AnswerS01E05{
		Task:   Task,
		ApiKey: ApiKey,
		Answer: Action,
	}

	messages := []openai.ChatCompletionMessageParamUnion{}
	var response http.Response

	systemPrompt := files.ReadFileToString(Task + "/" + "prompt.md")
	messages = append(messages, openai.SystemMessage(systemPrompt))

	// REPL
	iter := 0

	for {
		fmt.Printf("[+] ITERACJA %v\n", iter)
		if iter == maxIterations {
			break
		}

		responseOpenAI, err := openAiClient.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
			Model:    "gpt-5.4-mini",
			Tools:    tools,
			Messages: messages,
		})
		if err != nil {
			fmt.Println("Chat completion error.")
			fmt.Println(err)
			return
		}
		choice := responseOpenAI.Choices[0]

		messages = append(messages, openai.ChatCompletionMessageParamUnion{
			OfAssistant: &openai.ChatCompletionAssistantMessageParam{
				ToolCalls: llm.ToToolCallParams(choice.Message.ToolCalls),
			},
		})

		if len(choice.Message.ToolCalls) == 0 {
			fmt.Println("Tool calls equal zero error.")
			fmt.Println(choice.Message.Content)
			continue
		}

		for _, toolCall := range choice.Message.ToolCalls {
			fmt.Println("[+] Function " + toolCall.Function.Name + " chosen...")

			switch toolCall.Function.Name {
			case "help":
				var input types.HelpInput

				json.Unmarshal([]byte(toolCall.Function.Arguments), &input)
				input.Action = toolCall.Function.Name

				ans.Answer = input

			case "reconfigure":
				var input types.ReconfigureInput

				json.Unmarshal([]byte(toolCall.Function.Arguments), &input)
				input.Action = toolCall.Function.Name

				ans.Answer = input

			case "getstatus":
				var input types.GetStatusInput

				json.Unmarshal([]byte(toolCall.Function.Arguments), &input)
				input.Action = toolCall.Function.Name

				ans.Answer = input

			case "setstatus":
				var input types.SetStatusInput

				json.Unmarshal([]byte(toolCall.Function.Arguments), &input)
				input.Action = toolCall.Function.Name

				ans.Answer = input

			case "save":
				var input types.SetStatusInput

				json.Unmarshal([]byte(toolCall.Function.Arguments), &input)
				input.Action = toolCall.Function.Name

				ans.Answer = input

			case "success":
				var input types.SuccessInput
				json.Unmarshal([]byte(toolCall.Function.Arguments), &input)
				// input.Action = toolCall.Function.Name
				fmt.Println(input.Flag)

				return

			default:
				fmt.Println("Defalut...")

			}

			response = answer.SendAnswerReturnHttp(ans, Url_verify)

			body, err := io.ReadAll(response.Body)
			if err != nil {
				fmt.Printf("Error reading body: %v\n", err)
				return
			}

			b, _ := json.Marshal(body)
			result := string(b)
			// fmt.Println(result)
			var rawBase64 string
			json.Unmarshal(b, &rawBase64)
			resultBytes, _ := base64.StdEncoding.DecodeString(rawBase64)
			result = string(resultBytes)
			fmt.Printf("Result: %s\n", result)
			messages = append(messages, openai.ToolMessage(result, toolCall.ID))

			if response.StatusCode != 200 {
				if response.Header["Retry-After"] != nil {
					fmt.Printf("Retry-After: %v\n", response.Header["Retry-After"][0])
					penalty, _ := strconv.Atoi(response.Header["Retry-After"][0])
					waitingTime = penalty + 1
					fmt.Printf("Waiting %v seconds...\n", waitingTime)
					time.Sleep(time.Duration(waitingTime) * time.Second)
				} else { //if response.StatusCode == 503 {
					fmt.Printf("Waiting %v seconds...\n", waitingTime)
					time.Sleep(time.Duration(waitingTime) * time.Second)
				}

			}
		}
		time.Sleep(100 * time.Millisecond)
		fmt.Println("================================================================================")
		iter++
	}

}
