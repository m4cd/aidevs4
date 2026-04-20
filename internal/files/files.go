package files

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gocarina/gocsv"
)

func DownloadFile(path string, filename string, url string) (err error) {

	os.MkdirAll(path, os.FileMode(0755))

	filepath := path + "/" + filename

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func ReadFileToString(filename string) string {
	b, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return string(b)
}

// Unmarshalls CSV stored in a file in FilePath to an array
func UnmarshalCSV[T any](FilePath string) ([]*T, error) {

	csvFile, csvFileError := os.OpenFile(FilePath, os.O_RDWR, os.ModePerm)

	if csvFileError != nil {
		return nil, fmt.Errorf("failed to open CSV file: %w", csvFileError)
	}

	defer csvFile.Close()

	var result []*T

	if unmarshalError := gocsv.UnmarshalFile(csvFile, &result); unmarshalError != nil {
		return nil, fmt.Errorf("failed to open CSV file: %w", unmarshalError)
	}

	return result, nil
}

// Unmarshalls JSON stored in a file in FilePath
func UnmarshalJSON[T any](FilePath string) (*T, error) {
	jsonFile, jsonFileError := os.ReadFile(FilePath)

	if jsonFileError != nil {
		return nil, fmt.Errorf("failed to open CSV file: %w", jsonFileError)
	}

	var result *T

	unmarshalError := json.Unmarshal(jsonFile, &result)
	if unmarshalError != nil {
		return nil, fmt.Errorf("failed to open JSON file: %w", unmarshalError)
	}

	return result, nil
}
