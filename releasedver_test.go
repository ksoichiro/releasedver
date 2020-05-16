package releasedver

import (
	"flag"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, NewAnalyzer(), "a")
}

func TestAnalyzerWithPaths(t *testing.T) {
	testdata := analysistest.TestData()

	flags := flag.NewFlagSet("", flag.ExitOnError)
	flags.String("paths", "golang.org/x/tools", "")
	flags.String("root", "", "")
	analyzer := NewAnalyzer()
	analyzer.Flags = *flags

	analysistest.Run(t, testdata, analyzer, "b")
}

func TestAnalyzerWithMultiplePaths(t *testing.T) {
	testdata := analysistest.TestData()

	flags := flag.NewFlagSet("", flag.ExitOnError)
	flags.String("paths", "golang.org/x/tools,github.com/gostaticanalysis/modfile", "")
	flags.String("root", "", "")
	analyzer := NewAnalyzer()
	analyzer.Flags = *flags

	results := analysistest.Run(t, testdata, analyzer, "c")
	errors := 0
	for _, r := range results {
		errors += len(r.Diagnostics)
	}
	expectedErrors := 2
	if errors != expectedErrors {
		t.Errorf("got %d, want %d", errors, expectedErrors)
	}
}

func TestAnalyzerWithNoGoFilesInGoModDirectory(t *testing.T) {
	testdata := analysistest.TestData()

	flags := flag.NewFlagSet("", flag.ExitOnError)
	flags.String("paths", "golang.org/x/tools", "")
	flags.String("root", "", "")
	analyzer := NewAnalyzer()
	analyzer.Flags = *flags

	results := analysistest.Run(t, testdata, analyzer, "d/...")
	errors := 0
	for _, r := range results {
		errors += len(r.Diagnostics)
	}
	expectedErrors := 1
	if errors != expectedErrors {
		t.Errorf("got %d, want %d", errors, expectedErrors)
	}
}

func Test_parent(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{in: "a/b/c", want: "a/b"},
		{in: "a/b", want: "a"},
		{in: "a", want: ""},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.in, func(t *testing.T) {
			t.Parallel()
			s := parent(tt.in)
			if s != tt.want {
				t.Errorf("got %q, want %q", s, tt.want)
			}
		})
	}
}
