package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// categories maps a lowercase extension to the folder organize puts it in.
var categories = map[string]string{
	".jpg": "Images", ".jpeg": "Images", ".png": "Images", ".gif": "Images",
	".webp": "Images", ".heic": "Images", ".svg": "Images", ".bmp": "Images",

	".pdf": "PDFs",

	".doc": "Documents", ".docx": "Documents", ".txt": "Documents", ".md": "Documents",
	".rtf": "Documents", ".odt": "Documents", ".ppt": "Documents", ".pptx": "Documents",

	".csv": "Spreadsheets", ".xls": "Spreadsheets", ".xlsx": "Spreadsheets",

	".mp3": "Audio", ".wav": "Audio", ".m4a": "Audio", ".flac": "Audio",
	".mp4": "Video", ".mov": "Video", ".mkv": "Video", ".webm": "Video",
	".zip": "Archives", ".tar": "Archives", ".gz": "Archives", ".rar": "Archives", ".7z": "Archives",
	".exe": "Installers", ".msi": "Installers", ".dmg": "Installers", ".pkg": "Installers", ".deb": "Installers",
}

func category(name string) string {
	if c, ok := categories[strings.ToLower(filepath.Ext(name))]; ok {
		return c
	}
	return "Other"
}

// hidden reports whether a file should be left alone: dot files (the Unix
// convention) and the system files Windows keeps inside folders.
func hidden(name string) bool {
	return strings.HasPrefix(name, ".") ||
		strings.EqualFold(name, "desktop.ini") ||
		strings.EqualFold(name, "thumbs.db")
}

// dirArg returns args[i], or the current directory if it is missing.
func dirArg(args []string, i int) string {
	if i < len(args) {
		return args[i]
	}
	return "."
}


func humanSize(n int64) string {
	if n < 1024 {
		return fmt.Sprintf("%d B", n)
	}
	f, unit := float64(n)/1024, 0
	for f >= 1024 && unit < 4 {
		f /= 1024
		unit++
	}
	return fmt.Sprintf("%.1f %cB", f, "KMGTP"[unit])
}

// walkFiles calls fn for every regular, non-hidden file under root. also entries that
// cannot be read are reported and skipped rather than ending the walk.
func walkFiles(root string, fn func(path string, info fs.FileInfo) error) error {
	return filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			if path == root {
				return err
			}
			fmt.Fprintln(os.Stderr, "fm: skipping:", err)
			return nil
		}
		if path != root && hidden(d.Name()) {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if !d.Type().IsRegular() {
			return nil
		}
		info, err := d.Info()
		if err != nil {
			return err
		}
		return fn(path, info)
	})
}
