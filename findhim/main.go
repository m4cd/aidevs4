package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

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
	SuspectsCSVpath := "data/people/filteredPeople.csv"
	ApiKey := os.Getenv("API_KEY")
	Url_centrala := os.Getenv("URL_CNTRL")

	// Data
	PowerplantsFileName := "findhim_locations.json"
	PowerplantsFilePath := "data/findhim"
	PowerplantsJSONpath := PowerplantsFilePath + "/" + PowerplantsFileName
	Url_findhimLocation := Url_centrala + "data/" + ApiKey + "/" + PowerplantsFileName

	// APIs
	Url_location := Url_centrala + "api/location"
	Url_accesslevel := Url_centrala + "api/accesslevel"

	// REPL
	maxIterations := 10

	Url_verify := Url_centrala + "verify"

	OpenAiApiKey := os.Getenv("OPENAI_API_KEY")
	openAiClient := openai.NewClient(
		option.WithAPIKey(OpenAiApiKey),
	)

	err = files.DownloadFile(PowerplantsFilePath, PowerplantsFileName, Url_findhimLocation)
	if err != nil {
		fmt.Println("Error while downloading findhim JSON file.")
		return
	}

	// Suspects from S01E01
	suspects, err := files.UnmarshalCSV[types.Person](SuspectsCSVpath)
	if err != nil {
		return
	}

	// Prompt
	UserMessage := `Your goal: find which person is closest to any power plant, get their access level, and submit the answer.

To do this you need to:
- fetch all locations (coordinates) for every person
- fetch coordinates of every powerplant you're aware of
- calculate the nearest powerplant
- if the distance signifies the suspect was present at the powerplant then get the access level the suspect has
- identifiy the code of the powerplant (format PWR0000PL) and send the answer with the access level from previous step
- when in the response is a string in format {FLG:SOMETHING} present then call success function

<suspects>
`
	for _, s := range suspects {
		UserMessage = UserMessage + s.Name + " " + s.Surname + " born " + s.BirthDate + "\n"
	}
	UserMessage += "</suspects>\n\n"

	UserMessage += "<powerplants>\n"
	UserMessage += files.ReadFileToString(PowerplantsJSONpath) + "\n"
	UserMessage += "</powerplants>"

	messages := []openai.ChatCompletionMessageParamUnion{
		openai.UserMessage(UserMessage),
	}
	PrintMessages(messages)
	
	// REPL
	iter := 0
	for {
		fmt.Printf("[+] ITERACJA %v\n", iter)

		if iter == maxIterations {
			break
		}
		
		response, err := openAiClient.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
			Model:    "gpt-5",
			Tools:    tools,
			Messages: messages,
		})
		if err != nil {
			fmt.Println("Chat completion error.")
			fmt.Println(err)
			return
		}

		choice := response.Choices[0]

		messages = append(messages, openai.ChatCompletionMessageParamUnion{
			OfAssistant: &openai.ChatCompletionAssistantMessageParam{
				ToolCalls: toToolCallParams(choice.Message.ToolCalls),
			},
		})

		if len(choice.Message.ToolCalls) == 0 {
			fmt.Println("Tool calls equal zero error.")
			fmt.Println(choice.Message.Content)
			return
		}

		for _, toolCall := range choice.Message.ToolCalls {
			switch toolCall.Function.Name {
			case "get_powerplant_coordinates":
				fmt.Println("[+] Function \"get_powerplant_coordinates\" chosen...")
				var input types.ResolveCoordinatesInput
				json.Unmarshal([]byte(toolCall.Function.Arguments), &input)
				coords := resolveCoordinates(input.City)
				b, _ := json.Marshal(coords)
				result := string(b)
				fmt.Printf("Coordinates of %s powerplant: %s\n", input.City, result)
				messages = append(messages, openai.ToolMessage(result, toolCall.ID))

			case "nearest_powerplant":
				fmt.Println("[+] Function \"nearest_powerplant\" chosen...")
				var input nearestPowerplantInput
				json.Unmarshal([]byte(toolCall.Function.Arguments), &input)

				nearestpowerplant := nearestPowerplant(input)
				fmt.Printf("Nearest powerplant: %v\n", nearestpowerplant)
				b, _ := json.Marshal(map[string]float64{nearestpowerplant.city: nearestpowerplant.distance})
				messages = append(messages, openai.ToolMessage(string(b), toolCall.ID))

			case "get_suspects_locations":
				fmt.Println("[+] Function \"get_suspects_locations\" chosen...")

				var input LocationApiCallInput
				json.Unmarshal([]byte(toolCall.Function.Arguments), &input)
				fmt.Printf("Suspect: %s %s\n", input.Name, input.Surname)
				SuspectsLocations := LocationApiCall(ApiKey, Url_location, input.Name, input.Surname)

				b, err := json.Marshal(SuspectsLocations)
				if err != nil {
					fmt.Println("LocationApiCallInput marshal error.")
					return
				}
				messages = append(messages, openai.ToolMessage(string(b), toolCall.ID))

			case "get_accesslevel":
				fmt.Println("[+] Function \"get_accesslevel\" chosen...")

				var input AccessLevelApiCallInput

				json.Unmarshal([]byte(toolCall.Function.Arguments), &input)
				fmt.Printf("Suspect: %s %s\n", input.Name, input.Surname)
				AccessLevel := AccessLevelApiCall(ApiKey, Url_accesslevel, input.Name, input.Surname, input.Birthdate)

				b, err := json.Marshal(AccessLevel)
				if err != nil {
					fmt.Println("AccessLevelApiCallInput marshal error.")
					return
				}
				fmt.Printf("Suspects AccessLevel: %v\n", string(b))

				messages = append(messages, openai.ToolMessage(string(b), toolCall.ID))

			case "send_answer":
				fmt.Println("[+] Function \"send_answer\" chosen...")

				var input types.AnswerSuspectS01E02

				json.Unmarshal([]byte(toolCall.Function.Arguments), &input)
				fmt.Printf("Input: %v\n", input)
				fmt.Println(toolCall.Function.Arguments)
			
				ans := types.AnswerS01E02{
					Task:   "findhim",
					ApiKey: ApiKey,
					Answer: input,
				}
				fmt.Println(input)
				fmt.Println(ans)

				response := answer.SendAnswer(ans, Url_verify)
				fmt.Println(string(response))
				messages = append(messages, openai.ToolMessage(string(response), toolCall.ID))
			case "success":
				return
			default:
				fmt.Println("Defalut...")

			}

		}
		iter++
		fmt.Println()
	}

}

func PrintMessages(messages []openai.ChatCompletionMessageParamUnion) {
	for i, msg := range messages {
		switch {
		case msg.OfUser != nil:
			content := msg.OfUser.Content.OfString
			fmt.Printf("[%d] USER:\n%s\n\n", i, content)

		case msg.OfAssistant != nil:
			fmt.Printf("[%d] ASSISTANT:\n", i)
			if msg.OfAssistant.Content.OfString.Value != "" {
				fmt.Printf("  Text: %s\n", msg.OfAssistant.Content.OfString)
			}
			for _, tc := range msg.OfAssistant.ToolCalls {
				fmt.Printf("  Tool call: %s(%s) [id: %s]\n", tc.Function.Name, tc.Function.Arguments, tc.ID)
			}
			fmt.Println()

		case msg.OfTool != nil:
			fmt.Printf("[%d] TOOL RESULT [id: %s]:\n  %s\n\n", i, msg.OfTool.ToolCallID, msg.OfTool.Content.OfString)

		case msg.OfSystem != nil:
			fmt.Printf("[%d] SYSTEM:\n%s\n\n", i, msg.OfSystem.Content.OfString)
		}
	}
}

func toToolCallParams(toolCalls []openai.ChatCompletionMessageToolCall) []openai.ChatCompletionMessageToolCallParam {
	params := make([]openai.ChatCompletionMessageToolCallParam, len(toolCalls))
	for i, tc := range toolCalls {
		params[i] = openai.ChatCompletionMessageToolCallParam{
			ID:   tc.ID,
			Type: "function",
			Function: openai.ChatCompletionMessageToolCallFunctionParam{
				Name:      tc.Function.Name,
				Arguments: tc.Function.Arguments,
			},
		}
	}
	return params
}
