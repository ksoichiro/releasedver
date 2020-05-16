package main

import (
	"github.com/ksoichiro/releasedver"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(releasedver.Analyzer)
}
