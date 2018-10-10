package pjd

import "unicode"

func ToSnakeCase(in string) string {
	var out []rune
	out = append(out, unicode.ToLower(rune(in[0])))
	for i := 1; i < len(in); i++ {
		if rune(in[i]) == unicode.ToUpper(rune(in[i])) && rune(in[i-1]) != unicode.ToUpper(rune(in[i-1])) {
			out = append(out, rune('_'))
		} else if (i < len(in)-1) && rune(in[i+1]) != unicode.ToUpper(rune(in[i+1])) && rune(in[i]) == unicode.ToUpper(rune(in[i])) {
			out = append(out, rune('_'))
		}
		out = append(out, unicode.ToLower(rune(in[i])))
	}
	return string(out)
}
