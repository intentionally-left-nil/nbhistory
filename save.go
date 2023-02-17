package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func save(path string, commitMessage string) error {
	path, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	historyDir := path + ".history"
	err = createDirectory(historyDir)
	if err != nil {
		return err
	}
	err = ensureGit(historyDir)
	if err != nil {
		return err
	}
	notebook, err := openNotebook(path)
	if err != nil {
		return err
	}
	for _, cell := range notebook.Cells {
		err := saveCell(cell.Cell(), historyDir)
		if err != nil {
			return err
		}
	}
	err = saveNotebook(notebook, historyDir)
	if err != nil {
		return err
	}
	err = commitToGit(historyDir, commitMessage)
	if err != nil {
		return err
	}
	return nil
}

func createDirectory(path string) error {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		err = os.Mkdir(path, 0755)
		if err != nil {
			return fmt.Errorf("Could not create %s: %w", path, err)
		}
	}
	return err
}

func ensureGit(path string) error {
	_, err := os.Stat(filepath.Join(path, ".git"))
	if os.IsNotExist(err) {
		cmd := exec.Command("git", "init", path)
		err = cmd.Run()
	}
	return err
}

func stageFileInGit(path string) error {
	workDir := filepath.Dir(path)
	gitDir := filepath.Join(workDir, ".git")
	cmd := exec.Command("git", "--git-dir="+gitDir, "--work-tree="+workDir, "add", path)
	fmt.Printf("Calling %v\n", cmd)
	return cmd.Run()
}

func commitToGit(historyDir string, commitMessage string) error {
	gitDir := filepath.Join(historyDir, ".git")
	cmd := exec.Command("git", "--git-dir="+gitDir, "--work-tree="+historyDir, "commit", "-m", commitMessage)
	fmt.Printf("Calling %v\n", cmd)
	return cmd.Run()
}

func openNotebook(path string) (*Notebook, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	var notebook Notebook
	err = json.NewDecoder(file).Decode(&notebook)
	return &notebook, err
}

func saveCell(cell *Cell, path string) error {
	cellPath := filepath.Join(path, cell.Id)
	fmt.Printf("Saving to %s\n", cellPath)
	file, err := os.Create(cellPath)
	if err != nil {
		return err
	}
	defer file.Close()
	for _, line := range cell.Source.Data {
		_, err := file.Write([]byte(line))
		if err != nil {
			return err
		}
	}
	err = file.Close()
	if err != nil {
		return err
	}
	return stageFileInGit(cellPath)
}

func saveNotebook(notebook *Notebook, path string) error {
	notebookPath := filepath.Join(path, "notebook.ipynb")
	fmt.Printf("Saving to %s\n", notebookPath)
	file, err := os.Create(notebookPath)
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(notebook)
	if err != nil {
		return err
	}
	err = file.Close()
	if err != nil {
		return err
	}
	return stageFileInGit(notebookPath)
}
