package translation_data

import "fmt"

var data = map[string]map[string]string{
	"en": {
		"hello": "Hello",
	},
	"es": {
		"hello": "Hola",
	},
}

func GetTranslatedString(lang, text string) (string, error) {
	if result, ok := data[lang]; ok {
		if translated, ok := result[text]; ok {
			return translated, nil
		} else {
			return "", fmt.Errorf("Text not found: %s", text)
		}
	} else {
		return "", fmt.Errorf("Language not found: %s", lang)
	}
}
