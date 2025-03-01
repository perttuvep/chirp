package main

import (
	"strings"
)

func ChirpsValidate(str string) bool {

	const maxChirpLength = 140
	if len(str) > maxChirpLength {
		return false
	}
	return true

}

func profanity(s string) string {
	out := strings.Split(s, " ")
	for i, v := range out {
		if strings.ToLower(v) == "kerfuffle" || strings.ToLower(v) == "sharbert" || strings.ToLower(v) == "fornax" {
			out[i] = "****"
		}
	}
	return strings.Join(out, " ")
}
func ContainsAt[T comparable](list []T, item T) (bool, int) {
	for i, v := range list {
		if v == item {
			return true, i
		}
	}
	return false, -1
}
