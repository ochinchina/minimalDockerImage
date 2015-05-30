package main

import (
	"bufio"
	"filedep"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

func MakeTarball( tarName string, files []string ) {
	for index, file := range( files ) {
		if index == 0 {
			exec.Command( "tar", "-cvf", tarName, file ).Run()
		} else {
			exec.Command( "tar", "-rvf", tarName, file ).Run()
		}
	}
}

func LoadFiles( infiles string, fileName string )( *[]string, error) {

	f, err := os.Open( fileName )
	files := strings.Split( infiles, "," )
	if err != nil {
		return &files, nil
	}

	rd := bufio.NewReader( f )

	for {
		line, _, err := rd.ReadLine()
		if err != nil  {
			break
		}
		s := strings.TrimSpace( string( line ) )

		if len( s ) > 0 && s[0] != '#' {
			files = append( files, s )
		}
	}

	return &files, nil
}


func MakeImage( imageName string, files []string ) {
	f, err := ioutil.TempFile( ".", "mdi" )
	if err != nil {
		fmt.Printf( "error %v\n", err )
		return
	}

	f.Close()

	MakeTarball( f.Name(), files )

	catCmd := exec.Command( "cat", f.Name() )
	dockerImportCmd := exec.Command( "docker", "import", "-", imageName )

	//create pipe
	r, w := io.Pipe()

	dockerImportCmd.Stdin = r
	catCmd.Stdout = w

	catCmd.Start()
	dockerImportCmd.Start()

	catCmd.Wait()
	w.Close()

	dockerImportCmd.Wait()

	r.Close()

	os.Remove( f.Name() )
}


func printDepends( files []string ) {
	for _, file := range( files ) {
		fmt.Printf( "%s\n", file )
	}
}

func printUsage() {
	fmt.Printf( "Usage: mdi [-i <files>]  [-f <filename>] [-o <outputname>] command\n" )
	fmt.Printf( "files      files seperated by comma \",\"\n" )
	fmt.Printf( "filename   the absolute file name which contains all files should be packed\n")
	fmt.Printf( "outputname the output tar file name or docker image name\n")
	fmt.Printf( "command:\n")
	fmt.Printf( "    create_tar    create a tar file\n")
	fmt.Printf( "    create_image  create docker image\n")
	fmt.Printf( "    list_files    list all the files will ba packed\n" )
}

func parseFiles( files string )[]string {
	return strings.Split( files, "," )
}


func main() {

	var fileName = flag.String( "f", "", "file name" )
	var output = flag.String( "o", "", "output file name" )
	var files = flag.String( "i", "", "files separated by comma \",\"")
	flag.Parse()

	args := flag.Args()

	if len( args ) <= 0 {
		printUsage()
	} else {
		cmd := args[0]
		switch  cmd {
			case "create_tar":
				files, _ := LoadFiles( *files, *fileName )
				deps := filedep.FindDepend( files )
				MakeTarball( *output, deps )
			case "create_image":
				files, _ := LoadFiles( *files, *fileName )
				deps := filedep.FindDepend( files )
				MakeImage( *output, deps )
			case "list_files":
				files, _ := LoadFiles( *files, *fileName )
				deps := filedep.FindDepend( files )
				printDepends( deps )
			default:
				printUsage()
		}
	}
}

