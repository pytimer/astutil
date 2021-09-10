# astutil

Go ast util generate AST by parsing all go file of the project

## Usage

```go
package main

import (
	goparser "go/parser"
	"log"

	"github.com/pytimer/astutil"
)

func main() {
	pkgs := astutil.ParseDir("./testdata", []string{"vendor"}, goparser.ParseComments)
	for k, pkg := range pkgs {
		log.Printf("package path: %s, name: %s\n", k, pkg.Name)
		for _, f := range pkg.Files {
			log.Printf("astFile path: %s package: %s\n", f.Path, f.PackagePath)
			// log.Println(f.File)
			log.Printf("comments: %#v", f.Comments)
		}
	}
}
```