package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/joho/godotenv"
	"github.com/m4cd/aidevs4/internal/answer"
	"github.com/m4cd/aidevs4/internal/files"
	"github.com/m4cd/aidevs4/internal/maze"
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
	Task := "electricity"
	pngFilename := "electricity.png"
	pngSolvedFilename := "solved_electricity.png"

	pngSolvedURL := Url_centrala + "i/" + pngSolvedFilename
	pngURL := Url_centrala + "data/" + ApiKey + "/" + pngFilename

	pngPath := "data/electricity"

	fmt.Println("[+] Downloading current state PNG...")
	files.DownloadFile(pngPath, pngFilename, pngURL)

	fmt.Println("[+] Downloading desired state PNG...")
	files.DownloadFile(pngPath, pngSolvedFilename, pngSolvedURL)

	solutionMatrix, err := maze.SolveRotations(pngPath+"/"+pngFilename, pngPath+"/"+pngSolvedFilename)
	if err != nil {
		fmt.Println("SaveParts error" + err.Error())
		return
	}

	for r := 0; r < len(solutionMatrix); r++ {
		for c := 0; c < len(solutionMatrix[r]); c++ {
			fmt.Println(solutionMatrix)
			n := solutionMatrix[r][c]

			fmt.Printf("[+] Number of rotations of square %vx%v is  %v\n", r+1, c+1, n)
			for range n {

				rotate := types.AnswerRotateS02E02{
					Rotate: fmt.Sprintf("%vx%v", r+1, c+1),
				}
				ans := types.AnswerS02E02{
					Task:   Task,
					ApiKey: ApiKey,
					Answer: rotate,
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
		}
	}
}
