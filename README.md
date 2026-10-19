# go-cli-file-organizer

![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go&logoColor=white)
![Dependencies](https://img.shields.io/badge/dependencies-none-brightgreen)

A command-line tool that cleans up messy folders. It sorts files into folders by type, finds duplicate files even when their names are different, and can undo everything it moves.

I built it because my Downloads folder was always a mess, and as a way to learn Go.

## Features

- **Organize** files into folders like `Images`, `PDFs`, and `Documents`
- **Preview** every change before anything moves
- **Undo** everything with one command
- **Find duplicates** by content, not by name
- **Search** for files with patterns like `*.pdf`
- **See** what's using space, by file type

## Install

Requires Go 1.22 or newer.

```bash
git clone https://github.com/runweztt/go-cli-file-organizer
cd go-cli-file-organizer
go build -o fm .
```

On Windows, build with `go build -o fm.exe .` and run it as `.\fm.exe`.

## Commands

| Command | What it does |
|---|---|
| `fm organize [dir]` | Show where each file would go, without moving anything |
| `fm organize -apply [dir]` | Move the files into folders by type |
| `fm undo [dir]` | Put every moved file back |
| `fm stats [dir]` | Count files and space used per file type |
| `fm find "pattern" [dir]` | List files matching a pattern, like `"*.pdf"` |
| `fm dupes [dir]` | List files with identical content |

If you leave out the folder, `fm` uses the current one.

## Examples

### Organize a folder

```bash
fm organize ~/Downloads
```

```
bank_statement.pdf -> PDFs/bank_statement.pdf
budget_2026.csv -> Spreadsheets/budget_2026.csv
resume.docx -> Documents/resume.docx
vacation.jpg -> Images/vacation.jpg
4 files to move; run again with -apply to move them
```

Happy with the plan? Apply it:

```bash
fm organize -apply ~/Downloads
```

Changed your mind? Undo it:

```bash
fm undo ~/Downloads
```

### Find duplicates

```bash
fm dupes ~/Pictures
```

```
3 copies, 878.9 KB each:
  IMG_2041.jpg
  vacation (1).jpg
  vacation.jpg
```

All three files have different names, but the same content.

### See what's using space

```bash
fm stats ~/Downloads
```

```
TYPE          FILES  SIZE
Video         1      7.6 MB
Audio         1      3.3 MB
Images        4      2.9 MB
PDFs          2      361.3 KB
```

## Safety

- Nothing moves unless you add `-apply`.
- Existing files are never overwritten. If a file's destination is taken, it stays where it is.
- Every move is saved in `.fm-journal.json` inside the folder. Don't delete this file before running `undo`. Undo removes it when it's finished.
- Hidden files and Windows system files like `desktop.ini` are never touched.

## How it works

- **organize** plans every move before making any, so the preview shows exactly what `-apply` will do. Moving a file only renames it, so even large files move instantly.
- **dupes** groups files by size first, because files with different sizes can't match. It only reads and hashes (SHA-256) files that share a size, so most files are never opened.
- **stats**, **find**, and **dupes** share one folder walker that skips hidden folders like `.git`.

## Development

```bash
go vet ./...
go test ./...
```