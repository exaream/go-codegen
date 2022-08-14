package generator_test

import (
	"embed"
	"flag"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gostaticanalysis/skeletonkit"
	"github.com/tenntenn/golden"
)

const (
	modulePrefix = "example.com"
	outputDir    = "pkg"
	goldenDir    = "testdata/golden"
	tmplDir      = "testdata/template"
)

//go:embed testdata/template
var testTmplFS embed.FS

var (
	module     string
	flagUpdate bool
)

type pkgInfo struct {
	Name       string
	ModulePath string
}

func init() {
	flag.StringVar(&module, "module", "", "module")
	flag.BoolVar(&flagUpdate, "update", false, "update golden files")
}

// TODO: How to do force update.
// TODO: How to copy .gitkeep via skeletonkit.
// TODO: How to use *testing.T in main package.
func TestGenerator(t *testing.T) {
	if flagUpdate {
		golden.RemoveAll(t, goldenDir)
		os.RemoveAll(outputDir)
	}

	modulePath := path.Join(modulePrefix, module)
	moduleOutputDir := filepath.Join(outputDir, module)

	if err := os.MkdirAll(moduleOutputDir, 0755); err != nil {
		t.Fatal(err)
	}

	prompt := &skeletonkit.Prompt{
		Input:     strings.NewReader("a"),
		Output:    io.Discard,
		ErrOutput: io.Discard,
	}

	C := func(opts ...skeletonkit.CreatorOption) []skeletonkit.CreatorOption {
		return opts
	}
	creatorOpts := C(skeletonkit.CreatorWithPolicy(skeletonkit.ForceOverwrite))

	tmpl, err := skeletonkit.ParseTemplate(testTmplFS, module, tmplDir)
	if err != nil {
		t.Fatal(err)
	}

	fsys, err := skeletonkit.ExecuteTemplate(tmpl, pkgInfo{Name: module, ModulePath: modulePath})
	if err != nil {
		t.Fatal(err)
	}

	err = skeletonkit.CreateDir(prompt, moduleOutputDir, fsys, creatorOpts...)
	if err != nil {
		t.Fatal(err)
	}

	got := golden.Txtar(t, moduleOutputDir)

	if flagUpdate {
		golden.Update(t, goldenDir, module, got)
	}

	if diff := golden.Diff(t, goldenDir, module, got); diff != "" {
		t.Error(diff)
	}
}
