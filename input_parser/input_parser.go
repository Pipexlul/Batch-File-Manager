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
