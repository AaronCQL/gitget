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
			owner, name, found := strings.Cut(after, "/")
			if found {
				name, _ = strings.CutSuffix(name, ".git")
				return repo{
					owner: owner,
					name:  name,
				}, nil
			}
		}
	}
	return repo{}, fmt.Errorf("unable to parse URL: %v", s)
}
