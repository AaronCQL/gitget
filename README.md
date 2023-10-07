# `gitget`

A minimal CLI tool and library for downloading and unpacking [git archives](https://git-scm.com/docs/git-archive). Basically, the Go version of [`degit`](https://github.com/Rich-Harris/degit).

## Status

**Features**:

- [x] Automatically download the latest default branch (as opposed to assuming `main` or `master`)
- [x] Specify branches, tags, or commit hashes
- [ ] Specify subdirectories
- [ ] Support private repositories

**Providers**:

- [x] GitHub
- [ ] GitLab
- [ ] Bitbucket

## Installing

```sh
go install github.com/AaronCQL/gitget/cmd/gitget
```

## Usage

Downloads the latest default branch to the current working directory:

```sh
gitget https://github.com/owner/repo
```

The following are all equivalent:

```sh
gitget github:owner/repo
gitget git@github.com:owner/repo
gitget github.com/owner/repo
```

Specific branches, tags, and commit hashes can also be specified. Use the `--help` flag for more info:

```sh
gitget --help
```

## Go API

`gitget` can be used programmatically in your Go code. Start by installing this package:

```sh
go get -u github.com/AaronCQL/gitget
```

Then, use the `gitget.Clone` function:

```go
package main

import (
  "fmt"

  "github.com/AaronCQL/gitget/pkg/gitget"
)

func main() {
  res, err := gitget.Clone("github.com/owner/repo", gitget.Config{})
  if err != nil {
    panic(err)
  }
  fmt.Printf(
    "Cloned %v/%v (%v) into %v\n",
    res.RepoOwner, res.RepoName, res.RepoFragment, res.TargetDirRel,
  )
}
```
