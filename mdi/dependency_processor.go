package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
)

type DependencyProcessor interface {
	// process the dependencies in the list
	// returns error if any error occurs during the processing
	ProcessDependencies(deps *DependencyList) error
}

type DependencyTarballProcessor struct {
	outputFile string
}

func NewDependencyTarballProcessor(outputFile string) *DependencyTarballProcessor {
	return &DependencyTarballProcessor{outputFile: outputFile}
}

func (dtp *DependencyTarballProcessor) ProcessDependencies(deps *DependencyList) error {

	index := -1
	deps.ForEach(func(file string) {
		index++
		if index == 0 {
			exec.Command("tar", "-cf", dtp.outputFile, file).Run()
		} else {
			exec.Command("tar", "-uf", dtp.outputFile, file).Run()
		}
	})

	return nil
}

// create a temp tarball with the files in the DependencyList. The temp tarball should
// be deleted ASAP
//
// Returns:
//  the name of temp file and error indicator
func makeTempTarbar(deps *DependencyList) (string, error) {
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

type DependencyDockerImageMakeProcessor struct {
	imageName string
}

func NewDependencyDockerImageMakeProcessor(imageName string) *DependencyDockerImageMakeProcessor {
	return &DependencyDockerImageMakeProcessor{imageName: imageName}
}

func (ddp *DependencyDockerImageMakeProcessor) ProcessDependencies(deps *DependencyList) error {

	tarFileName, err := makeTempTarbar(deps)
	if err != nil {
		return err
	}

	defer os.Remove(tarFileName)

	return exec.Command("docker", "import", tarFileName, ddp.imageName).Run()
}

type DependencyPrintProcessor struct {
}

func (dpp DependencyPrintProcessor) ProcessDependencies(deps *DependencyList) error {
	deps.ForEach(func(dep string) {
		fmt.Println(dep)
	})
	return nil
}

type DependencyHttpOutProcessor struct {
	serverURL string
}

func (dhp DependencyHttpOutProcessor) ProcessDependencies(deps *DependencyList) error {
	tarFileName, err := makeTempTarbar(deps)
	if err != nil {
		return err
	}

	defer os.Remove(tarFileName)
	f, err := os.Open(tarFileName)

	if err != nil {
		return err
	}

	defer f.Close()

	resp, err := http.Post(fmt.Sprintf("%s/create_tar", tarFileName), "application/binary", f)
	if err != nil {
		return err
	}

	fmt.Println("Status:%s", resp.Status)
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	r := string(b)

	if len(r) > 0 {
		fmt.Println(r)
	}

	return nil
}
