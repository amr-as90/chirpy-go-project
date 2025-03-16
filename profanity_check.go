package main

import (
	"strings"
)

func ProfanityCheck(s string) string {
	splitString := strings.Split(s, " ")
	for splitWord := range splitString {
		lowerCasedWord := strings.ToLower(splitString[splitWord])
		if lowerCasedWord == "kerfuffle" || lowerCasedWord == "sharbert" || lowerCasedWord == "fornax" {
			splitString[splitWord] = "****"
		}
	}
	return strings.Join(splitString, " ")
}
