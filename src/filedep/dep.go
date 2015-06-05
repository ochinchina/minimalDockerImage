package filedep

import (
	"errors"
	"os/exec"
	"os"
	"path/filepath"
	"strings"
)

func FindDirectDepend( app string ) []string {
	r, err := exec.Command( "ldd", app ).Output()

	if err != nil {
		return []string{ app }
	}

	lines := strings.Split( string(r), "\n" )

	deps := make( []string, 0 )
	for _, line := range lines {
		fields := strings.Fields( line )
		dep := ""
		if len( fields  ) == 2 {
			if fields[1][0] == '(' {
				dep = fields[0]
			}
		} else if len( fields ) == 4 {
			dep = fields[2]
		}
		if dep != "" {
			deps = append( deps, dep )
		}

		for _, link := range FindLink( dep ) {
			deps = append( deps, link )
		}

	}
	return deps

}

func isDir( file string ) bool {
	fileInfo, err := os.Stat( file )
	if err != nil {
		return false
	}

	return fileInfo.IsDir()
}


func isExecutable( file string ) bool {
	fileInfo, err := os.Stat( file )
	 if err != nil {
                return false
        }

        return fileInfo.Mode() & 0111 != 0

}

func listFiles( path string ) []string {
	files := make( []string, 0 )
	filepath.Walk( path, func( file string, info os.FileInfo, err error ) error {
		if err != nil {
			return err
		}
		files = append( files, file )
		return nil
	})

	return files
}

func isUnderDir( path string, dir string ) bool {
	p1 := strings.Split( path, "/" )
	p2 := strings.Split( dir, "/" )

	if len( p1 ) >= len( p2 ) {
		for index, name := range p2 {
			if name != p1[index] {
				return false
			}
		}
		return true
	}
	return false
}

func FindDirDepend( dir string ) []string {
	deps := make( []string, 0 )
	for _, file := range listFiles( dir ) {
		if isExecutable( file ) {
			tmpDeps := FindDirectDepend( file )
			for _, dep := range tmpDeps {
				if !isUnderDir( dep, dir ) {
					deps = append( deps, dep )
				}
			}
		}
	}

	return deps
}


func addToDependMap( deps *[]string, depsMap *map[string]int, unProcessDeps *[]string ) {
	for _, dep := range *deps {
		_, ok := (*depsMap)[ dep ]

		if  !ok {
			*unProcessDeps = append( *unProcessDeps, dep )
			(*depsMap)[ dep ] = 1
		}
	}
}

func FindDepend( apps *[]string ) []string {
	depsMap := make( map[string]int )
	unProcessDeps := make( []string, 0)

	for _, app := range *apps {
		unProcessDeps = append( unProcessDeps, app )
		depsMap[ app ] = 1

		for _, link := range FindLink( app ) {
			depsMap[ link ] = 1
		}
	}

	for {
		if len( unProcessDeps ) == 0 {
			break
		}
		cur := unProcessDeps[0]
		unProcessDeps = unProcessDeps[1:]

		if isDir( cur ) {
			deps := FindDirDepend( cur )
			addToDependMap( &deps, &depsMap, &unProcessDeps )
		}  else if isExecutable( cur ) {
			deps := FindDirectDepend( cur )
			addToDependMap( &deps, &depsMap, &unProcessDeps )
		}
	}

	deps := make( []string, 0 )
	for k, _ := range depsMap {
		deps = append( deps, k )
	}

	return deps
}

func FindDirectLink( lib string )( string, error ) {
	r, err := exec.Command( "ls",  "-l", lib ).Output()

	if err != nil {
		return "", err
	}

	line := string( r )

	fields := strings.Fields( line )
	n := len( fields )

	if n > 3  && fields[ n - 2 ] == "->"  {

		if !filepath.IsAbs( fields[ n - 1 ] ) {

			i := strings.LastIndex( lib, "/" )

			if i > 0 {
				a := make( []string, 0 )
				a = append( a, lib[0:i] )
				a = append( a, fields[ n - 1 ]  )
				return strings.Join( a, "/" ), nil
			} 
		} else {
			return fields[ n - 1], nil
		}
	} 
	return "", errors.New( "no link found" )

}

func FindLink( lib string )[]string {
	cur := lib

	links := make( []string, 0 )

	for {
		link, err := FindDirectLink( cur )

		if err != nil {
			break	
		}
		if link != cur {
			links = append( links, link )
			cur = link
		} else {
			break
		}
	}

	return links
}

 
