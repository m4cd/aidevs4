package main

import (
	"encoding/base64"
	"fmt"
	"os"
)

func ImageFileToBase64(ImagePath string) string {

	_, err := os.Stat(ImagePath)
	if os.IsNotExist(err) {
		fmt.Println("IImagePath does not exist.")
		return ""
	}
	image, err := os.ReadFile(ImagePath)
	if err != nil {
		fmt.Println("Error while loading image file.")
		return ""
	}
	// os.WriteFile("base64.txt", []byte(base64.StdEncoding.EncodeToString(image)), 0644)

	imageBase64 := base64.StdEncoding.EncodeToString(image)

	return imageBase64
}
