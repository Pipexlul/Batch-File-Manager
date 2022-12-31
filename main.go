package main

import (
	"fmt"

	td "github.com/Pipexlul/batch-file-manager/translation_data"
)

func main() {
	fmt.Println("Testing translation data.")

	// Test English
	translated, err := td.GetTranslatedString("en", "hello")

}
