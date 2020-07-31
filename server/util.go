package server

import (
	"os"
	"path"
	"regexp"
	"strings"
)

func match(pattern string, b []byte) bool {
	re, err := regexp.Compile(pattern)
	if err != nil {
		log.Errorf("Error while math the path %s: %s", b, err)
		return false
	}

	return re.Match(b)
}

func resolvePath(p string) string {
	if strings.HasPrefix(p, "~/") {
		p = strings.TrimPrefix(p, "~/")

		home, err := os.UserHomeDir()
		if err != nil {
			log.Fatalf("Cannot get user home dir: ", err)
		}

		return path.Join(home, p)
	}

	return p
}
