package main

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/joho/godotenv"
	"github.com/m4cd/aidevs4/internal/answer"
	"github.com/m4cd/aidevs4/internal/files"
	"github.com/m4cd/aidevs4/internal/types"
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
	Task := "failure"
	logFilename := "failure.log"
	today := time.Now().Format("2006-01-02")
	prefixedFilename := today + "_" + logFilename

	logURL := Url_centrala + "data/" + ApiKey + "/" + logFilename
	logPath := "data/failure"

	fullPath := filepath.Join(logPath, prefixedFilename)
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		files.DownloadFile(logPath, prefixedFilename, logURL)
	}
	entries, err := parseLogFile(logPath + "/" + prefixedFilename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Błąd: %v\n", err)
		os.Exit(1)
	}

	var parsedEntries []LogEntry
	for _, e := range entries {
		if e.Level == "CRIT" {
			Exists := false
			for _, pe := range parsedEntries {

				if pe.Message == e.Message {

					Exists = true
				}

			}
			if !Exists {
				parsedEntries = append(parsedEntries, e)
			}

		}

	}

	logs := types.AnswerLogsS02E03{
		Logs: "",
	}
	for _, e := range parsedEntries {
		logs.Logs = logs.Logs + e.Raw + "\n"
	}
	fmt.Println(logs)

	ans := types.AnswerS02E03{
		ApiKey: ApiKey,
		Task:   Task,
		Answer: logs,
	}

	response := answer.SendAnswerReturnHttp(ans, Url_verify)

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

}

// LogEntry reprezentuje pojedynczy wpis logu
type LogEntry struct {
	Timestamp string
	Level     string
	Message   string
	Raw       string
}

// logRegex dopasowuje format: [TIMESTAMP] [LEVEL] MESSAGE
var logRegex = regexp.MustCompile(`^\[([^\]]+)\]\s+\[([^\]]+)\]\s+(.*)$`)

// parseLine parsuje pojedynczą linię logu
func parseLine(line string) (*LogEntry, bool) {
	matches := logRegex.FindStringSubmatch(line)
	if matches == nil {
		return nil, false
	}
	return &LogEntry{
		Timestamp: matches[1],
		Level:     matches[2],
		Message:   matches[3],
		Raw:       line,
	}, true
}

// parseLogFile czyta plik i zwraca wpisy inne niż INFO
func parseLogFile(path string) ([]LogEntry, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("nie można otworzyć pliku: %w", err)
	}
	defer file.Close()

	var entries []LogEntry
	scanner := bufio.NewScanner(file)
	// zwiększamy bufor na wypadek długich linii
	scanner.Buffer(make([]byte, 64*1024), 1024*1024)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		entry, ok := parseLine(line)
		if !ok {
			// linia nie pasuje do formatu - pomijamy
			continue
		}

		if entry.Level != "INFO" {
			entries = append(entries, *entry)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("błąd czytania pliku: %w", err)
	}

	return entries, nil
}
