package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

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
	Task := "drone"
	DroneMap := Url_centrala + "data/" + ApiKey + "/drone.png"


	OpenAiApiKey := os.Getenv("OPENAI_API_KEY")
	openAiClient := openai.NewClient(
		option.WithAPIKey(OpenAiApiKey),
	)
	params := openai.ChatCompletionNewParams{}
	var response http.Response

	systemPrompt := files.ReadFileToString(Task + "/" + "prompt.md")
	params.Messages = append(params.Messages, openai.SystemMessage(systemPrompt))

	systemPrompt = files.ReadFileToString(Task + "/" + "prompt_drone.md")
	params.Messages = append(params.Messages, openai.SystemMessage(systemPrompt))

	ChatCompletionContentPartImageParam := openai.ChatCompletionContentPartImageParam{
		ImageURL: openai.ChatCompletionContentPartImageImageURLParam{
			URL:    DroneMap,
			Detail: "high",
		},
		Type: "image_url",
	}
	params.Messages = append(params.Messages, openai.UserMessage([]openai.ChatCompletionContentPartUnionParam{
		{
			OfImageURL: &ChatCompletionContentPartImageParam,
		},
	}))

	params.Tools = tools
	params.Model = "gpt-5.4-mini"

	for {
		responseOpenAI, err := openAiClient.Chat.Completions.New(
			context.TODO(),
			params)
		if err != nil {
			fmt.Println("Chat completion error.")
			fmt.Println(err)
			return
		}
		choice := responseOpenAI.Choices[0]
		fmt.Println(choice.Message.Content)

		params.Messages = append(params.Messages, openai.ChatCompletionMessageParamUnion{
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
			fmt.Println(choice.Message.Content)

			switch toolCall.Function.Name {
			case "execute":
				var input types.AnswerInstructionsS02E05
				var ans types.AnswerS02E05
				ans.ApiKey = ApiKey
				ans.Task = Task

				json.Unmarshal([]byte(toolCall.Function.Arguments), &input)

				ans.Answer = input

				response = answer.SendAnswerReturnHttp(ans, Url_verify)

				result, _ := answer.PrintResponse(response)

				params.Messages = append(params.Messages, openai.ToolMessage(result, toolCall.ID))
			case "flg":
				type flag struct {
					Flg string `json:"flg"`
				}
				var input flag
				json.Unmarshal([]byte(toolCall.Function.Arguments), &input)

				structs.PrintStruct(input)
				return

			}
		}
	}
}
