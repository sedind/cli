// local build script file, similar to a makefile or collection of bash scripts in other projects
package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"github.com/sedind/cli"
)

func main() {
	app := cli.New()
	app.Description = "A local build script file, similar to a makefile or collection of bash scripts in other projects"

	app.Commands = cli.Commands{
		{
			Name:   "vet",
			Action: VetActionFunc,
		},
		{
			Name:   "test",
			Action: TestActionFunc,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

// VetActionFunc executes go vet command
func VetActionFunc(_ *cli.Context) error {
	return runCmd("go", "vet")
}

// TestActionFunc executes go test command
func TestActionFunc(c *cli.Context) error {

	packageName := "github.com/sedind/cli"

	coverProfile := "--coverprofile=coverprofile"

	err := runCmd("go", "test", "-v", coverProfile, packageName)
	if err != nil {
		return err
	}

	return testCleanup()
}

func testCleanup() error {
	var out bytes.Buffer

	file, err := os.Open("coverprofile")
	if err != nil {
		return err
	}

	b, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	out.Write(b)
	err = file.Close()
	if err != nil {
		return err
	}

	err = os.Remove("coverprofile")
	if err != nil {
		return err
	}

	outFile, err := os.Create("coverage.txt")
	if err != nil {
		return err
	}

	_, err = out.WriteTo(outFile)
	if err != nil {
		return err
	}

	err = outFile.Close()
	if err != nil {
		return err
	}

	return nil
}

func runCmd(arg string, args ...string) error {
	cmd := exec.Command(arg, args...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
