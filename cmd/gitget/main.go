package main

import (
	"fmt"
	"os"

	"github.com/AaronCQL/gitget/pkg/gitget"
	"github.com/spf13/pflag"
)

func main() {
	pflag.Usage = func() {
		fmt.Printf("Usage: gitget REPO_URL [OPTIONS...]\n")
		pflag.PrintDefaults()
	}

	cfg := gitget.Config{}
	pflag.StringVarP(&cfg.Dir, "dir", "d", "", "The target directory to clone into")
	pflag.StringVarP(&cfg.Branch, "branch", "b", "", "Git branch to clone")
	pflag.StringVarP(&cfg.Commit, "commit", "c", "", "Git commit hash to clone")
	pflag.StringVarP(&cfg.Tag, "tag", "t", "", "Git tag to clone")
	pflag.BoolVarP(&cfg.Force, "force", "f", false, "Forcefully write files into the existing target directory")
	pflag.BoolP("help", "h", false, "Display this help message")

	pflag.Parse()

	args := pflag.Args()
	if len(args) != 1 {
		pflag.Usage()
		os.Exit(1)
	}

	res, err := gitget.Clone(args[0], cfg)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Cloned %v/%v (%v) into %v\n", res.RepoOwner, res.RepoName, res.RepoFragment, res.TargetDirRel)
}
