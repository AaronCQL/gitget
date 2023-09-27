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

type Result struct {
	TargetDirRel string
	TargetDirAbs string
	RepoOwner    string
	RepoName     string
	RepoFragment string
}

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
		targetDir = filepath.Join(workDir, repo.name)
	} else if filepath.IsAbs(config.Dir) {
		targetDir = config.Dir
	} else {
		targetDir = filepath.Join(workDir, config.Dir)
	}

	// Check target dir
	if _, err := os.Stat(targetDir); err == nil && !config.Force {
		return Result{}, os.ErrExist
	}

	// Download tarball
	// Use the github api by default as it will fallback to the default branch
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
