package releasedver

import (
	"flag"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, Analyzer, "a")
}

func TestAnalyzerWithPaths(t *testing.T) {
	testdata := analysistest.TestData()

	flags := flag.NewFlagSet("", flag.ExitOnError)
	flags.String("paths", "golang.org/x/tools", "")
	analyzer := Analyzer
	analyzer.Flags = *flags

	analysistest.Run(t, testdata, analyzer, "b")
}

func TestAnalyzerWithMultiplePaths(t *testing.T) {
	testdata := analysistest.TestData()

	flags := flag.NewFlagSet("", flag.ExitOnError)
	flags.String("paths", "golang.org/x/tools,github.com/gostaticanalysis/modfile", "")
	analyzer := Analyzer
	analyzer.Flags = *flags

	analysistest.Run(t, testdata, analyzer, "c")
}
