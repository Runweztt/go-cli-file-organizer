package main

import "testing"

func TestCategory(t *testing.T) {
	tests := map[string]string{
		"photo.jpg":     "Images",
		"PHOTO.JPG":     "Images",
		"backup.tar.gz": "Archives",
		"Makefile":      "Other",
	}
	for name, want := range tests {
		if got := category(name); got != want {
			t.Errorf("category(%q) = %q, want %q", name, got, want)
		}
	}
}

func TestHumanSize(t *testing.T) {
	tests := map[int64]string{
		0:       "0 B",
		1023:    "1023 B",
		1024:    "1.0 KB",
		1536:    "1.5 KB",
		1 << 20: "1.0 MB",
		1 << 30: "1.0 GB",
	}
	for n, want := range tests {
		if got := humanSize(n); got != want {
			t.Errorf("humanSize(%d) = %q, want %q", n, got, want)
		}
	}
}
