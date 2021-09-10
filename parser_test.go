package astutil

import (
	goparser "go/parser"
	"testing"
)

func TestParseDir(t *testing.T) {
	pkgs := ParseDir("./testdata", []string{"vendor"}, goparser.ParseComments)
	for k, pkg := range pkgs {
		t.Logf("package path: %s, name: %s\n", k, pkg.Name)
		for _, f := range pkg.Files {
			t.Log("**************************************")
			t.Logf("astFile path: %s package: %s\n", f.Path, f.PackagePath)
			// t.Log(f.File)
			t.Logf("comments: %#v", f.Comments)
		}
	}
}
