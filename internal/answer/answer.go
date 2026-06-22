package answer

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func SendAnswer(answer any, URL string) []byte {
	httpClient := http.Client{}

	jsonBytes, err := json.Marshal(answer)
	if err != nil {
		fmt.Printf("Cannot marshal json: %s\n", err)
		os.Exit(1)
	}

	bodyReader := bytes.NewReader(jsonBytes)

	ansReq, err := http.NewRequest(http.MethodPost, URL, bodyReader)
	if err != nil {
		fmt.Printf("Cannot create response request: %s\n", err)
		os.Exit(1)
	}

	ansReq.Header.Set("Content-Type", "application/json")
	ansReq.Header.Set("User-Agent", "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:138.0) Gecko/20100101 Firefox/138.0")

	res, err := httpClient.Do(ansReq)

	if err != nil {
		fmt.Printf("Client error making http request: %s\n", err)
		os.Exit(1)
	}
	// defer res.Body.Close()

	if res.StatusCode != 200 {
		bodyBytes, _ := io.ReadAll(res.Body)
		fmt.Printf("Error %d: %s\n", res.StatusCode, string(bodyBytes)) // pokaże prawdziwy błąd
		os.Exit(1)
	}

	bodyAnswerBytes, _ := io.ReadAll(res.Body)

	return bodyAnswerBytes
}

func SendAnswerReturnHttp(answer any, URL string) http.Response {
	httpClient := http.Client{}

	jsonBytes, err := json.Marshal(answer)
	if err != nil {
		fmt.Printf("Cannot marshal json: %s\n", err)
		os.Exit(1)
	}

	bodyReader := bytes.NewReader(jsonBytes)

	ansReq, err := http.NewRequest(http.MethodPost, URL, bodyReader)
	if err != nil {
		fmt.Printf("Cannot create response request: %s\n", err)
		os.Exit(1)
	}

	ansReq.Header.Set("Content-Type", "application/json")
	ansReq.Header.Set("User-Agent", "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:138.0) Gecko/20100101 Firefox/138.0")

	res, err := httpClient.Do(ansReq)

	if err != nil {
		fmt.Printf("Client error making http request: %s\n", err)
		os.Exit(1)
	}
	// defer res.Body.Close()

	return *res
}

func SendPostJson[T any](endpoint string, payload *T) (string, error) {
	httpClient := http.Client{}

	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("Cannot marshal json: %s\n", err)
		return "", err
	}

	bodyReader := bytes.NewReader(jsonBytes)

	ansReq, err := http.NewRequest(http.MethodPost, endpoint, bodyReader)
	if err != nil {
		fmt.Printf("Cannot create response request: %s\n", err)
		return "", err
	}

	ansReq.Header.Set("Content-Type", "application/json")
	ansReq.Header.Set("User-Agent", "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:138.0) Gecko/20100101 Firefox/138.0")

	res, err := httpClient.Do(ansReq)

	if err != nil {
		fmt.Printf("Client error making http request: %s\n", err)
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		fmt.Printf("UnAuthorized with code: %d\n", res.StatusCode)
		return "", err
	}

	bodyAnswerBytes, _ := io.ReadAll(res.Body)

	return string(bodyAnswerBytes), nil

}

func PrintResponse(Response http.Response) (string, error) {
	body, err := io.ReadAll(Response.Body)
	if err != nil {
		fmt.Printf("Error reading body: %v\n", err)
		return "", err
	}
	b, _ := json.Marshal(body)
	result := string(b)

	var rawBase64 string
	json.Unmarshal(b, &rawBase64)
	resultBytes, _ := base64.StdEncoding.DecodeString(rawBase64)
	result = string(resultBytes)
	fmt.Printf("[+] Response: %s\n", result)
	return result ,nil
}
