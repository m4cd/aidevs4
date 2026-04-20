package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

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
	Task := "sendit"
	DocFileName := "index.md"
	DocURLlocation := Url_centrala + "dane/doc/"
	DocDataPath := "data/" + Task
	DocDataFile := DocDataPath + "/" + DocFileName
	AllowedImageExtensions := map[string]bool{
		".png": true,
	}
	AllowedTextExtensions := map[string]bool{
		".md": true,
	}

	// fmt.Println("===== variables ======")
	// fmt.Printf("Url_centrala: %s\n", Url_centrala)
	// fmt.Printf("Task: %s\n", Task)
	// fmt.Printf("DocFileName: %s\n", DocFileName)
	// fmt.Printf("DocURLlocation: %s\n", DocURLlocation)
	// fmt.Printf("DocURL: %s\n", DocURL)
	// fmt.Printf("DocDataPath: %s\n", DocDataPath)
	// fmt.Printf("DocDataFile: %s\n", DocDataFile)
	// fmt.Println("======================")

	OpenAiApiKey := os.Getenv("OPENAI_API_KEY")
	openAiClient := openai.NewClient(
		option.WithAPIKey(OpenAiApiKey),
	)
	ChatCompletionContentPartImageParam := openai.ChatCompletionContentPartImageParam{}

	// Documentation download
	DownloadNestedMarkdownFiles(DocDataPath, DocFileName, DocURLlocation)

	// Prompts definitions
	systemPrompt := files.ReadFileToString(Task + "/" + "prompt.md")

	userPrompt := "<dokumentacja>\n"

	// Adding main doc file:
	// DocFileName
	IndexMD, _ := os.ReadFile(DocDataFile)
	userPrompt = userPrompt + string(IndexMD)
	userPrompt = userPrompt + "</dokumentacja>\n\n"

	// Looping over data files
	Files, err := os.ReadDir(DocDataPath)
	if err != nil {
		fmt.Println("Cannot read DocDataPath directory.")
		fmt.Println(err)
		return
	}

	for _, f := range Files {
		fPath := filepath.Join(DocDataPath, f.Name())

		// index.md already included
		if fPath == DocDataFile {
			continue
		}

		userPrompt = userPrompt + fmt.Sprintf("<%v>\n", f.Name())

		ext := filepath.Ext(f.Name())
		if AllowedImageExtensions[ext] {
			// Image files
			userPrompt = userPrompt + "Obraz załączony.\n"
			ImagePath := DocDataPath + "/" + f.Name()
			Imgbase64 := ImageFileToBase64(ImagePath)
			ChatCompletionContentPartImageParam = openai.ChatCompletionContentPartImageParam{
				ImageURL: openai.ChatCompletionContentPartImageImageURLParam{
					URL:    fmt.Sprintf("data:image/jpeg;base64,%v", Imgbase64),
					Detail: "high",
				},
				Type: "image_url",
			}

		} else if AllowedTextExtensions[ext] {
			// Text files
			userPrompt = userPrompt + files.ReadFileToString(fPath)

		}
		userPrompt = userPrompt + fmt.Sprintf("</%v>\n\n", f.Name())
	}

	params := openai.ChatCompletionNewParams{}
	params.Messages = append(params.Messages, openai.SystemMessage(systemPrompt))
	params.Messages = append(params.Messages, openai.UserMessage(userPrompt))
	params.Messages = append(params.Messages, openai.UserMessage([]openai.ChatCompletionContentPartUnionParam{
		{
			OfImageURL: &ChatCompletionContentPartImageParam,
		},
	}))
	params.Model = openai.ChatModelGPT4o

	chatCompletion, err := openAiClient.Chat.Completions.New(
		context.TODO(),
		params,
	)
	if err != nil {
		fmt.Println("Chat completion error.")
		os.Exit(1)
	}

	Declaration := types.AnswerDeclarationS01E04{
		Declaration: chatCompletion.Choices[0].Message.Content,
	}

	ans := types.AnswerS01E04{
		Task:   "sendit",
		ApiKey: ApiKey,
		Answer: Declaration,
	}

	response := answer.SendAnswer(ans, Url_verify)
	fmt.Println(string(response))

}
