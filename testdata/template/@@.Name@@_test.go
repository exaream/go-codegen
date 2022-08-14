package sample_test

import (
	"testing"

	"@@.ModulePath@@"
)

func TestExplain(t *testing.T) {
	t.Parallel()

	want := "generated code for package @@.Name@@."
	got := @@.Name@@.Explain()

	if got != want {
		t.Errorf("want: %s, got: %s", want, got)
	}
}
