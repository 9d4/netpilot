package util

import "unicode"

func ToSnakeCase(input string) string {
	var output []rune
	for i, r := range input {
		if unicode.IsUpper(r) {
			if i > 0 && unicode.IsLower(rune(input[i-1])) {
				output = append(output, '_')
			}
			output = append(output, unicode.ToLower(r))
		} else {
			output = append(output, r)
		}
	}
	return string(output)
}
