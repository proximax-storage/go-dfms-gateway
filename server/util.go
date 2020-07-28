package server

import (
	"log"
	"regexp"
)

func match(pattern string, b []byte) bool {
	re, err := regexp.Compile(pattern)
	if err != nil {
		log.Println(err)
		return false
	}

	return re.Match(b)
}
