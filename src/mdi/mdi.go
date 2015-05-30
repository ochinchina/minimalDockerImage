package main

import (
	"bufio"
	"filedep"
	"flag"
	"fmt"
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

func LoadFiles( fileName string )( *[]string, error) {

	f, err := os.Open( fileName )
	if err != nil {
		return nil , err
	}

	files := make([]string, 0)

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


func printUsage() {
}

func MakeImage( imageName string, files []string ) {
}

func printDepends( files []string ) {
}
 
func main() {

	var fileName = flag.String( "f", "", "file name" )
	var output = flag.String( "o", "", "output file name" )
	flag.Parse()

	args := flag.Args()
	
	fmt.Printf( "fileName=%s\n", *fileName )
	fmt.Printf( "args=%v\n", args )

	if len( args ) <= 0 {
		printUsage()
	} else {
		cmd := args[0]

		switch  cmd {
			case "tar":
				files, _ := LoadFiles(  *fileName )	
				deps := filedep.FindDepend( files )
				MakeTarball( *output, deps )
			case "image":
				files, _ := LoadFiles(  *fileName ) 
				deps := filedep.FindDepend( files )
				MakeImage( *output, deps )
			case "list":
				files, _ := LoadFiles(  *fileName )
				deps := filedep.FindDepend( files )
				printDepends( deps )
		}
	}
}

