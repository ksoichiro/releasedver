package releasedver

import (
	"flag"
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"os"
	"regexp"
	"strings"

	"github.com/gostaticanalysis/modfile"
	xmodfile "golang.org/x/mod/modfile"
	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name:  "releasedver",
	Doc:   Doc,
	Run:   run,
	Flags: flags(),
	Requires: []*analysis.Analyzer{
		modfile.Analyzer,
	},
}

const Doc = "releasedver forces to use go modules with released version in go.mod file"

var releasedVerRE = regexp.MustCompile(`^v[^-]*$`)

func flags() flag.FlagSet {
	flags := flag.NewFlagSet("", flag.ExitOnError)
	flags.String("paths", "", "comma separated list of package paths to be checked")
	flags.String("root", "", "project root path, required when executed by go vet")
	return *flags
}

func run(pass *analysis.Pass) (interface{}, error) {
	paths := strings.Split(pass.Analyzer.Flags.Lookup("paths").Value.String(), ",")
	if len(paths) == 0 {
		return nil, nil
	}

	mf := pass.ResultOf[modfile.Analyzer].(*xmodfile.File)
	if mf == nil {
		// When modfile is not found, check the parent path to find go.mod.
		var dir string
		root := pass.Analyzer.Flags.Lookup("root").Value.String()
		if root != "" {
			// When executed by go vet and the directory that has go.mod has no Go files,
			// we cannot get current directory, so take it from flags.
			dir = root
		} else {
			var err error
			dir, err = os.Getwd()
			if err != nil {
				return nil, nil
			}
		}
		dir = strings.ReplaceAll(dir, "\\", "/")
		if !strings.Contains(dir, pass.Pkg.Path()) {
			p := types.NewPackage(parent(pass.Pkg.Path()), "a")
			parentPass := pass
			parentPass.Pkg = p

			result := map[*analysis.Analyzer]interface{}{}
			r, _ := modfile.Analyzer.Run(parentPass)
			result[modfile.Analyzer] = r
			parentPass.ResultOf = result
			run(parentPass)
		}
		return nil, nil
	}

	if len(mf.Require) > 0 {
		size := 0
		for _, r := range mf.Require {
			if size < r.Syntax.Start.Line {
				size = r.Syntax.Start.Line
			}
		}
		base := pass.Fset.Base()
		f := pass.Fset.AddFile(mf.Syntax.Name, pass.Fset.Base(), size)
		lines := make([]int, size)
		for i := range lines {
			lines[i] = i
		}
		f.SetLines(lines)

		comments := []*ast.CommentGroup{}
		for _, r := range mf.Require {
			if r.Syntax.Comments.Suffix != nil && len(r.Syntax.Comments.Suffix) > 0 {
				comments = append(comments, &ast.CommentGroup{
					List: []*ast.Comment{
						{
							Slash: token.Pos(base + r.Syntax.Start.Line - 1),
							Text:  r.Syntax.Comments.Suffix[0].Token,
						},
					},
				})
			}
		}
		file := &ast.File{
			Comments: comments,
		}
		pass.Files = append(pass.Files, file)

		for _, r := range mf.Require {
			if contains(paths, r.Mod.Path) && !releasedVerRE.MatchString(r.Mod.Version) {
				pass.Reportf(token.Pos(base+r.Syntax.Start.Line-1), fmt.Sprintf("%s must use released version", r.Mod.Path))
			}
		}
	}
	return nil, nil
}

func contains(array []string, s string) bool {
	for _, v := range array {
		if s == v {
			return true
		}
	}
	return false
}

func parent(path string) string {
	idx := strings.LastIndex(path, "/")
	if idx < 0 {
		return ""
	}
	return path[:idx]
}
