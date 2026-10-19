package main

import (
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
)

// runFind prints every file under dir whose name matches a shell glob.
// Matching ignores case, so "*.jpg" also finds "IMG_01.JPG".
func runFind(args []string) error {
	if len(args) == 0 {
		return errors.New(`find: missing pattern, e.g. fm find "*.pdf"`)
	}
	pattern := strings.ToLower(args[0])

	// Validate once so a typo is reported instead of silently matching nothing.
	if _, err := filepath.Match(pattern, ""); err != nil {
		return fmt.Errorf("find: bad pattern %q", args[0])
	}

	return walkFiles(dirArg(args, 1), func(path string, info fs.FileInfo) error {
		if ok, _ := filepath.Match(pattern, strings.ToLower(info.Name())); ok {
			fmt.Println(path)
		}
		return nil
	})
}
