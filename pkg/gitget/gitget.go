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

func Clone(repository string, options ...optionFn) (Result, error) {
	// Parse repo
	repo, err := parse(repository)
	if err != nil {
		return Result{}, err
	}

	// Init all opts
	opts := newOptions(options...)
	workDir, err := os.Getwd()
	if err != nil {
		return Result{}, err
	}
	targetDir := ""
	if opts.targetDir == "" {
		targetDir = filepath.Join(workDir, repo.name)
	} else if filepath.IsAbs(opts.targetDir) {
		targetDir = opts.targetDir
	} else {
		targetDir = filepath.Join(workDir, opts.targetDir)
	}

	// Check target dir
	if _, err := os.Stat(targetDir); err == nil && !opts.force {
		return Result{}, os.ErrExist
	}

	// Download tarball
	url := fmt.Sprintf("https://api.github.com/repos/%v/%v/tarball", repo.owner, repo.name)
	fragment := "HEAD"
	if opts.commit != "" {
		url = fmt.Sprintf("https://github.com/%v/%v/archive/%v.tar.gz", repo.owner, repo.name, opts.commit)
		fragment = opts.commit
	} else if opts.tag != "" {
		url = fmt.Sprintf("https://github.com/%v/%v/archive/refs/tags/%v.tar.gz", repo.owner, repo.name, opts.tag)
		fragment = opts.tag
	} else if opts.branch != "" {
		url = fmt.Sprintf("https://github.com/%v/%v/archive/refs/heads/%v.tar.gz", repo.owner, repo.name, opts.branch)
		fragment = opts.branch
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
			baseDir := strings.SplitN(header.Name, string(os.PathSeparator), 2)[1]
			if err := os.MkdirAll(filepath.Join(targetDir, baseDir), 0755); err != nil {
				return Result{}, err
			}
		case tar.TypeReg:
			baseDir := strings.SplitN(header.Name, string(os.PathSeparator), 2)[1]
			outFile, err := os.Create(filepath.Join(targetDir, baseDir))
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
