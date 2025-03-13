package string_utils

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

func ToSnakeCase(str string) string {
	var b strings.Builder
	var prev rune
	for i, v := range str {
		if unicode.IsLower(v) {
			b.WriteRune(v)
		} else {
			if i > 0 && (unicode.IsLower(prev) ||
				unicode.IsLower(nextRune(str[i+utf8.RuneLen(v):]))) {
				b.WriteByte('_')
			}
			b.WriteRune(unicode.ToLower(v))
		}
		prev = v
	}
	return b.String()
}

func nextRune(s string) rune { r, _ := utf8.DecodeRuneInString(s); return r }
