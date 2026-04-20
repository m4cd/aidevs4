package main

import (
	"fmt"
	"os"
	"regexp"

	"github.com/m4cd/aidevs4/internal/files"
)

func extractIncludes(content []byte) []string {
	re := regexp.MustCompile(`\[include\s+file="([^"]+)"\]`)
	matches := re.FindAllSubmatch(content, -1)

	var files []string
	for _, match := range matches {
		files = append(files, string(match[1]))
	}
	return files
}

func DownloadNestedMarkdownFiles(DocDataPath string, DocFileName string, DocURLlocation string) {
	DocDataFile := DocDataPath + "/" + DocFileName

	_, err := os.Stat(DocDataFile)
	if os.IsNotExist(err) {
		fmt.Println("Initial datafile does not exist. Downloading...")
		err = files.DownloadFile(DocDataPath, DocFileName, DocURLlocation+DocFileName)
		if err != nil {
			fmt.Println("DownloadFile error.")
			fmt.Println(err)
			return
		}
	}

	content, err := os.ReadFile(DocDataFile)
	if err != nil {
		fmt.Println("Reading DocDataFile error.")
		fmt.Println(err)
		return
	}

	// fmt.Println("[+] Loop over includes.")
	includes := extractIncludes(content)
	for _, f := range includes {
		url := DocURLlocation + f
		fPath := DocDataPath + "/" + f
		_, err = os.Stat(fPath)
		if os.IsNotExist(err) {
			fmt.Println("Data does not exist. Downloading...")
			err = files.DownloadFile(DocDataPath, f, url)
			if err != nil {
				fmt.Println("DownloadFile error in DownloadNestedMarkdownFiles.")
				fmt.Println(err)
				return
			}
		}
		DownloadNestedMarkdownFiles(DocDataPath, f, DocURLlocation)
	}
}
