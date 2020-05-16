package releasedver

import (
	"flag"
	"fmt"
	"go/ast"
	"go/token"
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
	return *flags
}

func run(pass *analysis.Pass) (interface{}, error) {
	paths := strings.Split(pass.Analyzer.Flags.Lookup("paths").Value.String(), ",")
	if len(paths) == 0 {
		return nil, nil
	}

	mf := pass.ResultOf[modfile.Analyzer].(*xmodfile.File)
	if mf != nil && len(mf.Require) > 0 {
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
