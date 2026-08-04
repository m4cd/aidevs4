package main

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/m4cd/aidevs4/internal/files"
	"github.com/m4cd/aidevs4/internal/types"
	"github.com/m4cd/aidevs4/internal/webserver"
	"github.com/openai/openai-go"
)

func Handler(key string, openAiClient openai.Client, mu *sync.RWMutex, Task string, ItemsCSVfile string, ConnectionsCSVfile string, CitiesCSVfile string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestBody := types.InputFromAgent{}

		err := json.NewDecoder(r.Body).Decode(&requestBody)
		if err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}

		messages := []openai.ChatCompletionMessageParamUnion{}
		systemPrompt := files.ReadFileToString(Task + "/" + "prompt.md")

		messages = append(messages, openai.SystemMessage(systemPrompt))
		messages = append(messages, openai.UserMessage(requestBody.Params))

		// llm.PrintAllMessages(messages)
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

		var ItemClasifier types.ItemS03E04
		if err := json.Unmarshal([]byte(llmtxt), &ItemClasifier); err != nil {
			fmt.Printf("parse error: %v", err)
		}

		// Patern matching for items
		f, err := os.Open(ItemsCSVfile)
		if err != nil {
			fmt.Fprintln(os.Stderr, "open:", err)
			os.Exit(1)
		}
		defer f.Close()

		csvReader := csv.NewReader(f)
		csvReader.FieldsPerRecord = -1 // tolerate rows with varying column counts

		lineNo := 0
		suspiciousItems := []string{}
		for {
			record, err := csvReader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				fmt.Fprintln(os.Stderr, "read:", err)
				os.Exit(1)
			}
			lineNo++

			// Match if any field contains the pattern
			for _, field := range record {
				if strings.Contains(field, ItemClasifier.Name) {
					line := strings.Join(record, ",")

					suspiciousItems = append(suspiciousItems, line)
				}
			}

		}

		for idx, si := range suspiciousItems {
			for _, techspec := range ItemClasifier.Parameters {
				if !strings.Contains(si, techspec.Value) {
					suspiciousItems = append(suspiciousItems[:idx], suspiciousItems[idx+1:]...)
				}
			}
		}

		if len(suspiciousItems) > 1 {
			fmt.Println("Too many suspicious items matching specs.")
			return
		}

		ItemCode := strings.Split(suspiciousItems[0], ",")[1]

		// Patern matching for connections
		f, err = os.Open(ConnectionsCSVfile)
		if err != nil {
			fmt.Fprintln(os.Stderr, "open:", err)
			os.Exit(1)
		}
		defer f.Close()

		csvReader = csv.NewReader(f)
		csvReader.FieldsPerRecord = -1 // tolerate rows with varying column counts

		lineNo = 0
		suspiciousCons := []string{}
		for {
			record, err := csvReader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				fmt.Fprintln(os.Stderr, "read:", err)
				os.Exit(1)
			}
			lineNo++

			// Match if any field contains the pattern
			for _, field := range record {
				if strings.Contains(field, ItemCode) {
					line := strings.Join(record, ",")

					suspiciousCons = append(suspiciousCons, line)
				}
			}

		}

		SuspiciousCityCodes := []string{}
		for _, suspCon := range suspiciousCons {
			SuspiciousCityCodes = append(SuspiciousCityCodes, strings.Split(suspCon, ",")[1])
		}

		// Patern matching for city names
		f, err = os.Open(CitiesCSVfile)
		if err != nil {
			fmt.Fprintln(os.Stderr, "open:", err)
			os.Exit(1)
		}
		defer f.Close()

		csvReader = csv.NewReader(f)
		csvReader.FieldsPerRecord = -1 // tolerate rows with varying column counts

		lineNo = 0
		CityNames := []string{}
		for {
			record, err := csvReader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				fmt.Fprintln(os.Stderr, "read:", err)
				os.Exit(1)
			}
			lineNo++

			// Match if any field contains the pattern
			for _, field := range record {
				for _, scc := range SuspiciousCityCodes {
					if strings.Contains(field, scc) {
						line := strings.Join(record, ",")

						CityNames = append(CityNames, strings.Split(line, ",")[0])
					}
				}

			}

		}

		ResponseToAgent := types.OutputToAgent{}
		ResponseToAgent.Output = strings.Join(CityNames, ",")

		webserver.RespondWithJSON(w, 200, ResponseToAgent)

	}
}

func HandlerFlag() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
