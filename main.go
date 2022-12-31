package main

import (
	"fmt"

	td "github.com/Pipexlul/batch-file-manager/translation_data"
)

func main() {
	fmt.Println("Testing translation data.")

	// Test English
	translated, err := td.GetTranslatedString("en", "hello")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(translated)
	}

	// Test Spanish
	translated, err = td.GetTranslatedString("es", "hello")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(translated)
	}
}
