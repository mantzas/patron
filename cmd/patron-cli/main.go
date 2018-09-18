package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

const (
	nameTemplate    = "{{name}}"
	versionTemplate = "{{version}}"
)

func main() {
	module := flag.String("m", "", `define the module name ("github.com/mantzas/patron")`)
	version := flag.String("v", "1.0.0", "define the version")
	path := flag.String("p", "", "define the project folder (defaults to current)")
	vendor := flag.Bool("d", true, "define vendoring behavior")
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

	name := nameFromModule(*module)
	log.Printf("assuming name: %s", name)

	err := createPathAndChdir(*path)
	if err != nil {
		log.Fatalf("failed to create path: %v", err)
	}

	err = setupGit()
	if err != nil {
		log.Fatalf("failed to create git: %v", err)
	}

	err = createMain(name, *version)
	if err != nil {
		log.Fatalf("failed to create main: %v", err)
	}

	err = createDockerfile(name)
	if err != nil {
		log.Fatalf("failed to create Dockerfile: %v", err)
	}

	err = createReadme(name)
	if err != nil {
		log.Fatalf("failed to create README.md: %v", err)
	}

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

func nameFromModule(module string) string {
	lst := strings.LastIndex(module, "/")
	if lst == -1 {
		return module
	}

	return module[lst+1:]
}

func setupGit() error {
	log.Printf("creating git repository")
	return exec.Command("git", "init").Run()
}

func createDockerfile(name string) error {
	return ioutil.WriteFile("Dockerfile", dockerfileContent(name), 0664)
}

func createReadme(name string) error {
	return ioutil.WriteFile("README.md", readmeContent(name), 0664)
}

func goMod(module string, vendor bool) error {
	log.Print("setup go module support")
	out, err := exec.Command("go", "mod", "init", module).CombinedOutput()
	log.Print(string(out))
	if err != nil {
		return errors.New(string(out))
	}
	if vendor {
		log.Print("add vendoring")
		out, err := exec.Command("go", "mod", "vendor").CombinedOutput()
		log.Print(string(out))
		if err != nil {
			return errors.New(string(out))
		}
	}
	return nil
}

func createPathAndChdir(path string) error {
	log.Printf("create folder: %s", path)
	err := os.MkdirAll(path, 0775)
	if err != nil {
		return err
	}

	log.Printf("cd into: %s", path)
	return os.Chdir(path)
}

func createMain(name, version string) error {
	folder := fmt.Sprintf("cmd/%s", name)
	log.Printf("create folder: %s", folder)
	err := os.MkdirAll(folder, 0775)
	if err != nil {
		return err
	}

	file := fmt.Sprintf("%s/main.go", folder)
	log.Printf("create file: %s", file)
	return ioutil.WriteFile(file, mainContent(name, version), 0664)
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

func dockerfileContent(name string) []byte {
	cnt := `FROM golang:latest as builder
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o {{name}} ./cmd/{{name}}/main.go 

FROM scratch
COPY --from=builder {{name}} .
CMD ["./{{name}}"]
`
	return []byte(strings.Replace(cnt, nameTemplate, name, -1))
}

func mainContent(name, version string) []byte {
	cnt := `package main

import (
	"fmt"
	"os"

	"github.com/mantzas/patron"
	"github.com/mantzas/patron/log"
)

func main() {
	name := "{{name}}"
	version := "{{version}}"

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
`
	cnt = strings.Replace(cnt, nameTemplate, name, -1)
	cnt = strings.Replace(cnt, versionTemplate, version, -1)
	return []byte(cnt)
}

func readmeContent(name string) []byte {
	cnt := `# {{name}}`
	return []byte(strings.Replace(cnt, nameTemplate, name, -1))
}
