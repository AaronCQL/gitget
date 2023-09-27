package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/AaronCQL/gitget/pkg/gitget"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: gitget REPO_URL\n")
		os.Exit(1)
	}

	var (
		repoUrl   string = os.Args[1]
		targetDir string
		force     bool
		branch    string
		commit    string
		tag       string
	)
	flag.StringVar(&targetDir, "dir", "", "The target directory to clone into")
	flag.BoolVar(&force, "force", false, "If the target directory already exists, forcefully write files into it")
	flag.StringVar(&branch, "branch", "", "Git branch to clone")
	flag.StringVar(&commit, "commit", "", "Git commit hash to clone")
	flag.StringVar(&tag, "tag", "", "Git tag to clone")

	if err := flag.CommandLine.Parse(os.Args[2:]); err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}

	res, err := gitget.Clone(repoUrl,
		gitget.WithTargetDir(targetDir),
		gitget.WithForce(force),
		gitget.WithBranch(branch),
		gitget.WithCommit(commit),
		gitget.WithTag(tag),
	)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Cloned %v/%v (%v) into %v\n", res.RepoOwner, res.RepoName, res.RepoFragment, res.TargetDirRel)
}
