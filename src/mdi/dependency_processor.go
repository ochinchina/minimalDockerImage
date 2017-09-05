package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

type DependencyProcessor interface {
	ProcessDependencies(deps []string) error
}

type DependencyTarballProcessor struct {
	outputFile string
}

func NewDependencyTarballProcessor(outputFile string) *DependencyTarballProcessor {
	return &DependencyTarballProcessor{outputFile: outputFile}
}

func (dtp *DependencyTarballProcessor) ProcessDependencies(deps []string) error {

	for index, file := range deps {
		if index == 0 {
			exec.Command("tar", "-cvf", dtp.outputFile, file).Run()
		} else {
			exec.Command("tar", "-rvf", dtp.outputFile, file).Run()
		}
	}
	return nil
}

type DependencyDockerImageMakeProcessor struct {
	imageName string
}

func NewDependencyDockerImageMakeProcessor(imageName string) *DependencyDockerImageMakeProcessor {
	return &DependencyDockerImageMakeProcessor{imageName: imageName}
}

func (ddp *DependencyDockerImageMakeProcessor) makeTarbar(deps []string) (string, error) {
	f, err := ioutil.TempFile(".", "mdi")
	if err != nil {
		return "", err
	}

	f.Close()

	dtp := NewDependencyTarballProcessor(f.Name())
	err = dtp.ProcessDependencies(deps)
	if err != nil {
		os.Remove(f.Name())
		return "", err
	}
	return f.Name(), nil
}
func (ddp *DependencyDockerImageMakeProcessor) ProcessDependencies(deps []string) error {

	tarFileName, err := ddp.makeTarbar(deps)
	if err != nil {
		return err
	}

	defer os.Remove(tarFileName)

	return exec.Command("docker", "import", tarFileName).Run()
}

type DependencyPrintProcessor struct {
}

func (dpp DependencyPrintProcessor) ProcessDependencies(deps []string) error {
	fmt.Println(strings.Join(deps, "\n"))
	return nil
}
