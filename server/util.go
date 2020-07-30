package server

import (
	"log"
	"os"
	"path"
	"regexp"
	"strings"
)

func match(pattern string, b []byte) bool {
	re, err := regexp.Compile(pattern)
	if err != nil {
		log.Println(err)
		return false
	}

	return re.Match(b)
}

func resolvePath(p string) string {
	if strings.HasPrefix(p, "~/") {
		p = strings.TrimPrefix(p, "~/")

		home, err := os.UserHomeDir()
		if err != nil {
			log.Fatal("Cannot get user home dir: ", err)
		}

		return path.Join(home, p)
	}

	return p
}
