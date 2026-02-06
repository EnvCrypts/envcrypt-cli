package cmd

import (
	"fmt"
	"os/exec"
	"strings"
)



func getRepoFromGit() (string, error) {
	out, err := exec.Command("git", "remote", "get-url", "origin").Output()
	if err != nil {
		return "", err
	}

	url := strings.TrimSpace(string(out))
	url = strings.TrimSuffix(url, ".git")
	url = strings.ReplaceAll(url, "git@github.com:", "")
	url = strings.ReplaceAll(url, "https://github.com/", "")

	return url, nil
}

func getCurrentBranch() (string, error) {
	out, err := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD").Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func buildRepoPrincipal(repo, branch string) string {
	return fmt.Sprintf(
		"repo:%s:ref:refs/heads/%s",
		repo,
		branch,
	)
}

func DetectGitContext() (string, string, string, error) {
	repo, err := getRepoFromGit()
	if err != nil {
		return "", "", "", fmt.Errorf("could not detect git repo")
	}

	branch, err := getCurrentBranch()
	if err != nil {
		return "", "", "", fmt.Errorf("could not detect git branch")
	}

	principal := buildRepoPrincipal(repo, branch)
	
	return principal, repo, branch, nil
}
