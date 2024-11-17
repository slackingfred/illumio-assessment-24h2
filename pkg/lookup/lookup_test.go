package lookup

import (
	"path"
	"runtime"
	"testing"
)

func TestSampleFile(t *testing.T) {
	{
		tbl := Table{}
		_, testSourceFile, _, _ := runtime.Caller(0)
		if err := tbl.LoadFile(path.Join(path.Dir(testSourceFile), "../../testdata/lookup_table.csv")); err != nil {
			t.Errorf("Error loading file: %v", err)
		}
		if tag := tbl.GetTag(25, "Tcp"); tag != "sv_P1" {
			t.Errorf("Expected `sv_P1`, got %v", tag)
		}
		if tag := tbl.GetTag(80, "tcp"); tag != "" {
			t.Errorf("Expected empty, got %v", tag)
		}
	}
}
