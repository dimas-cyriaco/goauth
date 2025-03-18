package string_utils

import "github.com/leonklingele/randomstring"

func GenerateSecureRandomString(length int) (string, error) {
	return randomstring.Generate(length, randomstring.CharsAlphaNum)
}
