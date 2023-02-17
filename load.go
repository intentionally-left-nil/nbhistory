package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

func revert(path string) error {
	path, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	historyDir := path + ".history"
	gitDir := filepath.Join(historyDir, ".git")
	cmd := exec.Command("git", "--git-dir="+gitDir, "--work-tree="+historyDir, "reset", "--hard", "HEAD~1")
	fmt.Printf("Calling %v\n", cmd)
	return cmd.Run()
}

func load(path string) error {
	path, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	historyDir := path + ".history"
	backupPath := filepath.Join(historyDir, "notebook.ipynb")
	fmt.Printf("Restoring from %s\n", backupPath)
	src, err := os.Open(backupPath)
	if err != nil {
		return err
	}
	defer src.Close()

	fmt.Printf("Restoring to %s\n", path)
	dest, err := os.Create(path)
	if err != nil {
		return err
	}
	defer dest.Close()
	_, err = io.Copy(dest, src)
	return err
}
