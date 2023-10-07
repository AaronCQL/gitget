package gitget

import (
	"fmt"
	"strings"
)

var prefixes = []string{
	"github:",
	"git@github.com:",
	"https://github.com/",
	"github.com/",
}

type repo struct {
	owner string
	name  string
}

func parse(s string) (repo, error) {
	for _, prefix := range prefixes {
		after, ok := strings.CutPrefix(s, prefix)
		if ok {
			split := strings.Split(after, "/")
			if len(split) >= 2 {
				name, _ := strings.CutSuffix(split[1], ".git")
				return repo{
					owner: split[0],
					name:  name,
				}, nil
			}
		}
	}
	return repo{}, fmt.Errorf("unable to parse URL: %v", s)
}
