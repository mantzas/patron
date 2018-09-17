package main

import (
	"errors"
	"log"
	"os"
	"os/exec"
)

func main() {
	err := setupGit()
	if err != nil {
		log.Fatalf("failed to create git: %v", err)
	}

	err = createStructure()
	if err != nil {
		log.Fatalf("failed to create folder structure: %v", err)
	}

	// Create cmd and main.go with initial code

	// Create docker file

	// Create dashboard folder and files

	// create readme.md

	err = setupGoMod()
	if err != nil {
		log.Fatalf("failed to initialize go mod support: %v", err)
	}

	err = initialGitCommit()
	if err != nil {
		log.Fatalf("failed to commit initially to git: %v", err)
	}
}

func setupGit() error {
	log.Printf("creating git...")
	return exec.Command("git", "init").Run()
}

func setupGoMod() error {
	log.Print("initializing go module support...")
	out, err := exec.Command("go", "mod", "init").CombinedOutput()
	if err != nil {
		return errors.New(string(out))
	}
	return nil
}

func createStructure() error {
	return os.MkdirAll("cmd", 0775)
}

func initialGitCommit() error {
	log.Printf("initial git commit...")
	return exec.Command("git", "commit", "-m", "Initial commit").Run()
}
