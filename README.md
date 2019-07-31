# git-private-org-repo-ls

Command line utility to list the repositories available within a private organization's scope on Github

## Getting Started

* Download the builds tarball/zip from the [releases page](https://github.com/komish/git-private-org-repo-ls/releases).
* Extract tarball/zip and find the relevant binary for your OS-architecture.
* Place at a preferred location in your user path and mark executable (if necessary).
* Generate a [Github Personal Access Token](https://github.com/settings/tokens) with the following permissions:
    * repo:status
    * repo_deployment
    * public_repo
    * repo:invite
    * read:org 
* Place this token at `${HOME}/.gittoken` (required). 

NOTE: An improperly scoped token will return a successful return code (because the token is accepted) but will return no repository links or data.

## Usage


```
# git-private-org-repo-ls --help
List repositories that exist within a private organization on Github.

Requires a personal access token be placed in $HOME/.gittoken
	https://github.com/settings/tokens

Required token scope:
	repo:status, repo_deployment, public_repo, repo:invite, read:org

Usage: ./git-private-org-repo-ls [OPTIONS] orgname
  -L	Print your Github API limits.
  -v	Display the version of this utility.
```

## Example Output

```
# git-private-org-repo-ls supercool-org
-----------------------------------
         Name: supercool-repo
         Link: https://github.com/supercool-org/supercool-repo
  Clone (SSH): git@github.com:supercool-org/supercool-repo
Clone (HTTPS): https://github.com/supercool-org/supercool-repo
-----------------------------------
(...)
```

## Build for Multiple Platforms

Binaries are built using the following process for multiple platforms.

* Clone this repository.
* Change directory to repository path.
* Use `go get` to build external dependencies.
* Run build commands to compile for various platforms.

```
env GOOS=linux GOARCH=amd64 go build -o builds/linux-amd64/git-private-org-repo-ls
env GOOS=linux GOARCH=386 go build -o builds/linux-386/git-private-org-repo-ls
env GOOS=linux GOARCH=arm go build -o builds/linux-arm/git-private-org-repo-ls
env GOOS=linux GOARCH=arm64 go build -o builds/linux-arm64/git-private-org-repo-ls
env GOOS=windows GOARCH=amd64 go build -o builds/windows-amd64/git-private-org-repo-ls.exe
env GOOS=darwin GOARCH=amd64 go build -o builds/darwin-amd64/git-private-org-repo-ls
```
