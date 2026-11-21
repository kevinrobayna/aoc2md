package internal

import (
	"flag"
	"os"
	"path/filepath"
	"testing"
)

var update = flag.Bool("update", false, "update golden files")

func TestConvertDescription(t *testing.T) {
	tests := []struct {
		name    string
		fixture string
		golden  string
	}{
		{name: "both parts", fixture: "description.html", golden: "description.golden.md"},
		{name: "no description", fixture: "empty.html", golden: "empty.golden.md"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in, err := os.Open(filepath.Join("testdata", tt.fixture))
			if err != nil {
				t.Fatal(err)
			}
			defer in.Close()

			got, err := convertDescription(in)
			if err != nil {
				t.Fatalf("convertDescription() error = %v", err)
			}

			goldenPath := filepath.Join("testdata", tt.golden)
			if *update {
				if err := os.WriteFile(goldenPath, []byte(got), 0o600); err != nil {
					t.Fatal(err)
				}
			}
			want, err := os.ReadFile(goldenPath)
			if err != nil {
				t.Fatal(err)
			}
			if got != string(want) {
				t.Errorf("convertDescription() mismatch\n--- got ---\n%s\n--- want ---\n%s", got, want)
			}
		})
	}
}
