package translation_data

import "fmt"

var staticData = map[string]map[string]string{
	"en": {
		"hello":      "Hello",
		"langChosen": "Language chosen: English",
	},
	"es": {
		"hello":      "Hola",
		"langChosen": "Idioma elegido: Español",
	},
}

var data map[string]map[string]string

func init() {
	data = make(map[string]map[string]string)

	for lang, text := range staticData {
		data[lang] = make(map[string]string)
		for key, value := range text {
			data[lang][key] = value
		}
	}
	// TODO: Add command data
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
