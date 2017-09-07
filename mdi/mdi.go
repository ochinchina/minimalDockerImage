package main

import (
	"flag"
	"fmt"
	"strings"
)

func LoadConfig(infiles string, fileName string) *ImageConfig {
	config, err := NewImageConfigWithFile(fileName)
	if err != nil {
		config.includes = strings.Split(infiles, ",")
	}
	return config
}

func printUsage() {
	fmt.Printf("Usage: mdi [-i <files>]  [-f <filename>] [-o <outputname>] command\n")
	fmt.Printf("files      files seperated by comma \",\"\n")
	fmt.Printf("filename   the absolute file name which contains all files should be packed\n")
	fmt.Printf("outputname the output tar file name or docker image name\n")
	fmt.Printf("command:\n")
	fmt.Printf("    create_tar    create a tar file, the option -o must be present\n")
	fmt.Printf("    create_image  create docker image, the option -o will be ignored if present\n")
	fmt.Printf("    list_files    list all the files will ba packed, the option -o will be ignored if present\n")
	fmt.Printf("\n")
	fmt.Printf("For the option -i and -f, you can only choose one, don't input both\n")
}

func main() {

	var fileName = flag.String("f", "", "file name")
	var output = flag.String("o", "", "output file name")
	var files = flag.String("i", "", "files separated by comma \",\"")
	flag.Parse()

	args := flag.Args()

	if len(args) <= 0 {
		printUsage()
	} else {
		cmd := args[0]
		config := LoadConfig(*files, *fileName)
		dependencyFinder := NewDependencyFinder(config)
		deps := dependencyFinder.FindDependencies()
		processors := map[string]DependencyProcessor{"create_tar": NewDependencyTarballProcessor(*output),
			"create_image": NewDependencyDockerImageMakeProcessor(*output),
			"list_files":   &DependencyPrintProcessor{}}

		if processor, ok := processors[cmd]; ok {
			processor.ProcessDependencies(deps)
		} else {
			printUsage()
		}
	}
}
