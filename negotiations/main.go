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

	Task := "negotiations"

	ItemsCSVfile := "data/" + Task + "/" + "items.csv"
	ConnectionsCSVfile := "data/" + Task + "/" + "connections.csv"
	CitiesCSVfile := "data/" + Task + "/" + "cities.csv"
	

	OpenAiApiKey := os.Getenv("OPENAI_API_KEY")
	openAiClient := openai.NewClient(
		option.WithAPIKey(OpenAiApiKey),
	)

	mu := &sync.RWMutex{}

	serverPort := "8080"

	server := webserver.CreateWebserver(map[string]http.HandlerFunc{
		"/api": Handler(ApiKey, openAiClient, mu, Task, ItemsCSVfile, ConnectionsCSVfile, CitiesCSVfile),
		"/flagapi": HandlerFlag(),
	}, serverPort)

	server.ListenAndServe()

}
