package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/m4cd/aidevs4/internal/answer"
	"github.com/m4cd/aidevs4/internal/files"
	"github.com/m4cd/aidevs4/internal/llm"
	"github.com/m4cd/aidevs4/internal/structs"
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
	Task := "mailbox"
	MailApi := Url_centrala + "api/zmail"

	OpenAiApiKey := os.Getenv("OPENAI_API_KEY")
	openAiClient := openai.NewClient(
		option.WithAPIKey(OpenAiApiKey),
	)

	messages := []openai.ChatCompletionMessageParamUnion{}
	var response http.Response

	systemPrompt := files.ReadFileToString(Task + "/" + "prompt.md")
	messages = append(messages, openai.SystemMessage(systemPrompt))

	for {
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

		if len(choice.Message.ToolCalls) == 0 && len(choice.Message.ToolCalls) > 1 {
			fmt.Println("Tool calls different than 1 error.")
			fmt.Println(choice.Message.Content)
			continue
		}

		for _, toolCall := range choice.Message.ToolCalls {
			fmt.Println("[+] Function " + toolCall.Function.Name + " chosen...")

			switch toolCall.Function.Name {
			case "help":
				var input types.MailboxHelp
				input.Page = 1
				input.ApiKey = ApiKey

				json.Unmarshal([]byte(toolCall.Function.Arguments), &input)
				input.Action = toolCall.Function.Name

				response = answer.SendAnswerReturnHttp(input, MailApi)
				result, _ := answer.PrintResponse(response)
				messages = append(messages, openai.ToolMessage(result, toolCall.ID))

			case "getInbox":
				var input types.MailboxGetInbox
				input.Page = 1
				input.PerPage = 5
				input.ApiKey = ApiKey

				json.Unmarshal([]byte(toolCall.Function.Arguments), &input)
				input.Action = toolCall.Function.Name

				response = answer.SendAnswerReturnHttp(input, MailApi)
				result, _ := answer.PrintResponse(response)
				messages = append(messages, openai.ToolMessage(result, toolCall.ID))

			case "search":
				var input types.MailboxSearch
				input.Page = 1
				input.PerPage = 5
				input.ApiKey = ApiKey

				json.Unmarshal([]byte(toolCall.Function.Arguments), &input)
				input.Action = toolCall.Function.Name

				response = answer.SendAnswerReturnHttp(input, MailApi)
				result, _ := answer.PrintResponse(response)
				messages = append(messages, openai.ToolMessage(result, toolCall.ID))

			case "getMessages":
				var input types.MailboxGetMessages
				input.ApiKey = ApiKey

				json.Unmarshal([]byte(toolCall.Function.Arguments), &input)
				input.Action = toolCall.Function.Name

				response = answer.SendAnswerReturnHttp(input, MailApi)
				result, _ := answer.PrintResponse(response)
				messages = append(messages, openai.ToolMessage(result, toolCall.ID))

			case "getThread":
				var input types.MailboxGetThread
				input.ApiKey = ApiKey

				json.Unmarshal([]byte(toolCall.Function.Arguments), &input)
				input.Action = toolCall.Function.Name

				response = answer.SendAnswerReturnHttp(input, MailApi)
				result, _ := answer.PrintResponse(response)
				messages = append(messages, openai.ToolMessage(result, toolCall.ID))

			case "reset":
				var input types.MailboxHelp
				input.ApiKey = ApiKey

				json.Unmarshal([]byte(toolCall.Function.Arguments), &input)
				input.Action = toolCall.Function.Name

				response = answer.SendAnswerReturnHttp(input, MailApi)
				result, _ := answer.PrintResponse(response)
				messages = append(messages, openai.ToolMessage(result, toolCall.ID))

			case "verify":
				var input types.AnswerMessageS02E04
				var ans types.AnswerS02E04
				ans.ApiKey = ApiKey
				ans.Task = Task

				json.Unmarshal([]byte(toolCall.Function.Arguments), &input)

				structs.PrintStruct(input)

				ans.Answer = input

				response = answer.SendAnswerReturnHttp(ans, Url_verify)
				answer.PrintResponse(response)
				return
			}
		}
		time.Sleep(1 * time.Second)
	}

}
