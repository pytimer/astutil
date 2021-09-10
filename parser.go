package astutil

import (
	"fmt"
	"go/ast"
	goparser "go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

type parser struct {
	excludes  []string
	pkgs      map[string]*PackageDefinition
	parseMode goparser.Mode
}

type AstFileInfo struct {
	// Path the path of the ast.File
	Path string
	File *ast.File
	// PackagePath package import path of the ast.File
	PackagePath string
	Comments    ast.CommentMap
}

type PackageDefinition struct {
	// package name
	Name  string
	Files map[string]*AstFileInfo
}

func ParseDir(dir string, excludeDirs []string, mode goparser.Mode) map[string]*PackageDefinition {
	p := &parser{
		pkgs:      make(map[string]*PackageDefinition),
		parseMode: mode,
		excludes:  excludeDirs,
	}
	// search dir all go file
	filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {

		if err := p.skip(dir, path, f); err != nil {
			return err
		}
		if f.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}

		packageDir := filepath.ToSlash(filepath.Dir(filepath.Clean(filepath.Join(dir, relPath))))
		p.parseFile(packageDir, path, nil)

		return nil
	})

	return p.pkgs
}

func (p *parser) parseFile(packageDir, path string, src interface{}) error {
	fset := token.NewFileSet()
	astFile, err := goparser.ParseFile(fset, path, src, p.parseMode)
	if err != nil {
		return err
	}
	// ast.Print(fset, astFile)

	var comments ast.CommentMap
	if p.parseMode == goparser.ParseComments {
		comments = ast.NewCommentMap(fset, astFile, astFile.Comments)
	}

	if p.pkgs == nil {
		p.pkgs = make(map[string]*PackageDefinition)
	}

	path, err = filepath.Abs(path)
	if err != nil {
		return err
	}

	pk, ok := p.pkgs[packageDir]
	if ok {
		if _, exists := pk.Files[path]; exists {
			return nil
		}
		pk.Files[path] = &AstFileInfo{
			Path:        path,
			File:        astFile,
			PackagePath: packageDir,
			Comments:    comments,
		}
	} else {
		p.pkgs[packageDir] = &PackageDefinition{
			Name: astFile.Name.Name,
			Files: map[string]*AstFileInfo{
				path: {
					Path:        path,
					File:        astFile,
					PackagePath: packageDir,
					Comments:    comments,
				},
			},
		}
	}

	return nil
}

func ParseAstComments(commentMap ast.CommentMap) map[string][]*ast.Comment {
	comments := make(map[string][]*ast.Comment)
	for node := range commentMap {
		var name string
		switch n := node.(type) {
		case *ast.FuncDecl:
			name = n.Name.Name
			if n.Recv != nil {
				if len(n.Recv.List) > 0 {
					recv := n.Recv.List[0]
					switch pointer := recv.Type.(type) {
					case *ast.StarExpr:
						if x, ok := pointer.X.(*ast.Ident); ok {
							name = fmt.Sprintf("(%s *%s) %s", recv.Names[0].Name, x.Name, name)
						}

					case *ast.Ident:
						name = fmt.Sprintf("(%s %s) %s", recv.Names[0].Name, pointer.Name, name)
					}
				}
			}

		case *ast.GenDecl:
			// If n.Doc is no-nil, Specs only have one item,
			// otherwise should loop specs to get every var/const comment.
			valueSpec := n.Specs[0].(*ast.ValueSpec)
			name = valueSpec.Names[0].Name
		case *ast.ValueSpec:
			name = n.Names[0].Name
		}
		cs := make([]*ast.Comment, 0)
		for _, comment := range commentMap.Filter(node).Comments() {
			cs = append(cs, comment.List...)
		}
		comments[name] = cs
	}
	return comments
}

func (p *parser) skip(dir, path string, f os.FileInfo) error {
	if f.IsDir() {
		if f.Name() == "vendor" || (len(f.Name()) > 1 && f.Name()[0] == '.') {
			return filepath.SkipDir
		}
		rel, err := filepath.Rel(dir, path)
		if err != nil {
			return filepath.SkipDir
		}
		if p.excludes != nil && inSlice(rel, p.excludes) {
			return filepath.SkipDir
		}
	} else {
		if filepath.Ext(f.Name()) != ".go" || strings.HasSuffix(strings.ToLower(f.Name()), "_test.go") {
			return filepath.SkipDir
		}
	}
	return nil
}

func inSlice(s string, slice []string) bool {
	for _, str := range slice {
		if s == str {
			return true
		}
	}
	return false
}
