package gitget

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// The configuration for the `Clone` function.
type Config struct {
	// The target directory to clone into. If empty, the repository name is used.
	Dir string
	// The Git branch to clone.
	Branch string
	// The Git tag to clone.
	Tag string
	// The Git commit hash to clone.
	Commit string
	// Forcefully write files into the existing target directory.
	Force bool
}

// The result of the `Clone` function.
type Result struct {
	// The relative path to the target directory.
	TargetDirRel string
	// The absolute path to the target directory.
	TargetDirAbs string
	// The owner of the repository, typically the user or organisation.
	RepoOwner string
	// The name of the repository.
	RepoName string
	// The commit, tag, or branch used to clone the repository.
	RepoFragment string
}

// Clones the given repository using the configuration provided. The default behaviour
// clones the HEAD of the default branch.
//
//	res, err := gitget.Clone("github.com/AaronCQL/gitget", gitget.Config{})
//	if err != nil {
//		panic(err)
//	}
//	fmt.Printf(
//		"Cloned %v/%v (%v) into %v\n",
//		res.RepoOwner, res.RepoName, res.RepoFragment, res.TargetDirRel,
//	)
func Clone(repository string, config Config) (Result, error) {
	// Parse repo string
	repo, err := parse(repository)
	if err != nil {
		return Result{}, err
	}

	// Init all options
	workDir, err := os.Getwd()
	if err != nil {
		return Result{}, err
	}
	targetDir := ""
	if config.Dir == "" {
		// Use current working dir and repo name as default dir name
		targetDir = filepath.Join(workDir, repo.name)
	} else if filepath.IsAbs(config.Dir) {
		// If the given dir is absolute, use it as is
		targetDir = config.Dir
	} else {
		// If the given dir is relative, use it as relative to the working dir
		targetDir = filepath.Join(workDir, config.Dir)
	}

	// Check target dir
	if _, err := os.Stat(targetDir); err == nil && !config.Force {
		return Result{}, fmt.Errorf("target directory already exists: %v", targetDir)
	}

	// Form the archive URL
	// Use the github api by default as it uses the default branch
	url := fmt.Sprintf("https://api.github.com/repos/%v/%v/tarball", repo.owner, repo.name)
	fragment := "HEAD"
	if config.Commit != "" {
		url = fmt.Sprintf("https://github.com/%v/%v/archive/%v.tar.gz", repo.owner, repo.name, config.Commit)
		fragment = config.Commit
	} else if config.Tag != "" {
		url = fmt.Sprintf("https://github.com/%v/%v/archive/refs/tags/%v.tar.gz", repo.owner, repo.name, config.Tag)
		fragment = config.Tag
	} else if config.Branch != "" {
		url = fmt.Sprintf("https://github.com/%v/%v/archive/refs/heads/%v.tar.gz", repo.owner, repo.name, config.Branch)
		fragment = config.Branch
	}

	// Download the tarball
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return Result{}, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return Result{}, err
	}
	if res.StatusCode == 404 {
		return Result{}, fmt.Errorf("respository not found: %v", repository)
	}
	if res.StatusCode >= 400 {
		return Result{}, fmt.Errorf("server replied with status code: %v", res.StatusCode)
	}

	// Unzip and untar the tarball
	gzipReader, err := gzip.NewReader(res.Body)
	if err != nil {
		return Result{}, err
	}
	tarReader := tar.NewReader(gzipReader)
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return Result{}, err
		}
		switch header.Typeflag {
		case tar.TypeDir:
			targetPath := strings.SplitN(header.Name, string(os.PathSeparator), 2)[1]
			if err := os.MkdirAll(filepath.Join(targetDir, targetPath), 0755); err != nil {
				return Result{}, err
			}
		case tar.TypeReg:
			targetPath := strings.SplitN(header.Name, string(os.PathSeparator), 2)[1]
			outFile, err := os.Create(filepath.Join(targetDir, targetPath))
			if err != nil {
				return Result{}, err
			}
			if _, err := io.Copy(outFile, tarReader); err != nil {
				outFile.Close()
				return Result{}, err
			}
			if err := outFile.Close(); err != nil {
				return Result{}, err
			}
		case tar.TypeXGlobalHeader:
			// ignore these headers that are in git archives
		default:
			return Result{}, fmt.Errorf("unsupported header %v in tar file entry %v", header.Typeflag, header.Name)
		}
	}

	targetDirRel, _ := strings.CutPrefix(targetDir, workDir)
	return Result{
		TargetDirRel: targetDirRel,
		TargetDirAbs: targetDir,
		RepoOwner:    repo.owner,
		RepoName:     repo.name,
		RepoFragment: fragment,
	}, nil
}
