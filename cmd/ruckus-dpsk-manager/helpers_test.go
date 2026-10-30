package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolvePassword(t *testing.T) {
	dir := t.TempDir()
	write := func(name, content string) string {
		p := filepath.Join(dir, name)
		if err := os.WriteFile(p, []byte(content), 0o600); err != nil {
			t.Fatal(err)
		}
		return p
	}
	withNewline := write("pw", "s3cret\n")
	crlf := write("crlf", "s3cret\r\n")
	empty := write("empty", "\n")

	cases := []struct {
		name, flag, file, want string
		wantErr                bool
	}{
		{"flag", "s3cret", "", "s3cret", false},
		{"file trims newline", "", withNewline, "s3cret", false},
		{"file trims crlf", "", crlf, "s3cret", false},
		{"both", "x", withNewline, "", true},
		{"neither", "", "", "", true},
		{"empty file", "", empty, "", true},
		{"missing file", "", filepath.Join(dir, "nope"), "", true},
	}
	for _, c := range cases {
		got, err := resolvePassword(c.flag, c.file)
		if (err != nil) != c.wantErr || got != c.want {
			t.Errorf("%s: got %q, %v", c.name, got, err)
		}
	}
}
