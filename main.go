package main

import (
	"fmt"

	"github.com/Pipexlul/batch-file-manager/translation_data"
)

func main() {
	fmt.Println("Testing translation data.")

	// Test English
	translated, err := translation_data.GetTranslatedString("en", "hello")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(translated)
	}

	// Test Spanish
	translated, err = translation_data.GetTranslatedString("es", "hello")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(translated)
	}
}
