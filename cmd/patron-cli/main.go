package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

func main() {
	module := flag.String("m", "", `define the module name ("github.com/mantzas/patron")`)
	path := flag.String("p", "", "define the project folder (defaults to current)")
	vendor := flag.Bool("v", true, "define vendoring behavior")
	flag.Parse()

	if *path == "" {
		fmt.Print("path is required\n\n")
		flag.Usage()
		os.Exit(1)
	}

	if *module == "" {
		fmt.Print("module is required\n\n")
		flag.Usage()
		os.Exit(1)
	}

	err := createPathAndChdir(*path)
	if err != nil {
		log.Fatalf("failed to create path: %v", err)
	}

	err = setupGit()
	if err != nil {
		log.Fatalf("failed to create git: %v", err)
	}

	err = createMain()
	if err != nil {
		log.Fatalf("failed to create main: %v", err)
	}

	// Create docker file

	// Create dashboard folder and files

	// create readme.md

	err = goMod(*module, *vendor)
	if err != nil {
		log.Fatalf("failed to initialize go mod support: %v", err)
	}

	err = goFormat()
	if err != nil {
		log.Fatalf("failed to execute go fmt: %v", err)
	}

	err = gitCommit()
	if err != nil {
		log.Fatalf("failed to commit initially to git: %v", err)
	}
	log.Print("completed successful")
}

func setupGit() error {
	log.Printf("creating git repository")
	return exec.Command("git", "init").Run()
}

func goMod(module string, vendor bool) error {
	log.Print("setup go module support")
	out, err := exec.Command("go", "mod", "init", module).CombinedOutput()
	if err != nil {
		return errors.New(string(out))
	}
	if vendor {
		log.Print("add vendoring")
		out, err := exec.Command("go", "mod", "vendor").CombinedOutput()
		if err != nil {
			return errors.New(string(out))
		}
	}
	return nil
}

func createPathAndChdir(path string) error {
	log.Printf("create %s folder", path)
	err := os.MkdirAll(path, 0775)
	if err != nil {
		return err
	}

	log.Printf("cd to %s", path)
	return os.Chdir(path)
}

func createMain() error {
	log.Print("create cmd folder")
	err := os.MkdirAll("cmd", 0775)
	if err != nil {
		return err
	}

	log.Print("create main.go")
	return ioutil.WriteFile("cmd/main.go", mainContent(), 0664)
}

func gitCommit() error {
	log.Printf("git add")
	err := exec.Command("git", "add", ".").Run()
	if err != nil {
		return err
	}
	log.Printf("git commit")
	return exec.Command("git", "commit", "-m", "Initial commit").Run()
}

func goFormat() error {
	log.Print("running go fmt")
	out, err := exec.Command("go", "fmt", "./...").CombinedOutput()
	if err != nil {
		return errors.New(string(out))
	}
	return nil
}

func mainContent() []byte {
	return []byte(`package main

import (
	"fmt"
	"os"

	"github.com/mantzas/patron"
	"github.com/mantzas/patron/log"
)

func main() {
	name := "patron"
	version := "1.0.0"

	err := patron.SetupLogging(name, version)
	if err != nil {
		fmt.Printf("failed to set up logging: %v", err)
		os.Exit(1)
	}

	srv, err := patron.New(name, version)
	if err != nil {
		log.Fatalf("failed to create service %v", err)
	}

	err = srv.Run()
	if err != nil {
		log.Fatalf("failed to run service %v", err)
	}
}
`)
}
