package input_parser

import "strings"

func ParseInput(input string) map[string]string {
	args := make(map[string]string)

	tokens := strings.Fields(input)
	tokensLen := len(tokens)

	for i := 0; i < tokensLen; i++ {
		if strings.HasPrefix(tokens[i], "-") {
			argName := strings.TrimPrefix(tokens[i], "-")

			if i+1 < tokensLen && !strings.HasPrefix(tokens[i+1], "-") {
				argValue := tokens[i+1]

				args[argName] = argValue
				i++
			} else {
				args[argName] = ""
			}
		}
	}
	return args
}

func ParseInputNew(input string) map[string]string {
	args := make(map[string]string)

	tokens := make([]string, 0)

	for i := 0; i < len(input); i++ {
		if input[i] == '"' {
			i++
			var token string

			for i < len(input) && input[i] != '"' {
				token += string(input[i])
				i++
			}

			tokens = append(tokens, token)
		} else if input[i] == ' ' {
			continue
		} else {
			var token string

			for i < len(input) && input[i] != ' ' {
				token += string(input[i])
				i++
			}

			tokens = append(tokens, token)
		}
	}

	tokensLen := len(tokens)

	for i := 0; i < tokensLen; i++ {
		if strings.HasPrefix(tokens[i], "-") {
			argName := strings.TrimPrefix(tokens[i], "-")

			if i+1 < tokensLen && !strings.HasPrefix(tokens[i+1], "-") {
				argValue := tokens[i+1]

				args[argName] = argValue
				i++
			} else {
				args[argName] = ""
			}
		}
	}
	return args
}
