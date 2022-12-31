package translation_data

import "testing"

func TestGetTranslatedString(t *testing.T) {
	t.Run("Ensure both languages have the same amount of keys", func(t *testing.T) {
		if len(data["en"]) != len(data["es"]) {
			t.Errorf("English and Spanish don't have the same amount of keys")
		}
	})

	t.Run("Ensure basic keys are being translated correctly and check errors", func(t *testing.T) {
		testData := []struct {
			lang     string
			text     string
			expected string
		}{
			{"en", "hello", "Hello"},    // Test English
			{"es", "hello", "Hola"},     // Test Spanish
			{"fr", "hello", ""},         // Test non-existent language
			{"en", "randomStuffGo", ""}, // Test non-existent text
		}

		for _, test := range testData {
			translated, err := GetTranslatedString(test.lang, test.text)
			if err != nil {
				if test.expected != "" {
					t.Errorf("Error: %s", err)
				}
			} else {
				if translated != test.expected {
					t.Errorf("Expected %s, got %s", test.expected, translated)
				}
			}
		}
	})
}
