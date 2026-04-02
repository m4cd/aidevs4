package main

import (
	"fmt"
	"net/http"
	"os"
	"sync"

	"github.com/joho/godotenv"
	"github.com/m4cd/aidevs4/internal/webserver"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

/*
*/

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error while loading .env file.")
		return
	}

	// // Variables
	ApiKey := os.Getenv("API_KEY")
	Url_centrala := os.Getenv("URL_CNTRL")
	maxIterations := 5

	// APIs
	Url_packages := Url_centrala + "api/packages"


	OpenAiApiKey := os.Getenv("OPENAI_API_KEY")
	openAiClient := openai.NewClient(
		option.WithAPIKey(OpenAiApiKey),
	)

	sessions := make(map[string]*Session)
	mu := &sync.RWMutex{}

	serverPort := "8080"
	
	server := webserver.CreateWebserver(map[string]http.HandlerFunc{
		"/api":    Handler(sessions, ApiKey, Url_packages, openAiClient, maxIterations, mu),
	
	}, serverPort)

	server.ListenAndServe()

}
