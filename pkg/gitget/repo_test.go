package gitget

import "testing"

func TestParse(t *testing.T) {
	validStrings := []string{
		"github:AaronCQL/gitget",
		"github:AaronCQL/gitget.git",
		"git@github.com:AaronCQL/gitget",
		"git@github.com:AaronCQL/gitget.git",
		"https://github.com/AaronCQL/gitget",
		"https://github.com/AaronCQL/gitget.git",
		"github.com/AaronCQL/gitget",
		"github.com/AaronCQL/gitget.git",
		"https://github.com/AaronCQL/gitget/blob/main/something/random.git",
		"https://github.com/AaronCQL/gitget.git/blob/main/something/random.git",
	}

	for _, s := range validStrings {
		parsed, err := parse(s)
		if err != nil {
			t.Error(err)
		}
		if parsed.owner != "AaronCQL" {
			t.Error("owner mismatch")
		}
		if parsed.name != "gitget" {
			t.Error("name mismatch")
		}
	}
}
