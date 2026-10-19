go-cli-file-organizer

A small command-line tool that sorts a messy folder into folders by file type, finds duplicate files by their content, and can undo everything it moves.

I built it because my Downloads folder was always a mess, and as a way to learn Go. It uses only the standard library.

Install
bash
git clone https://github.com/runweztt/go-cli-file-organizer
cd go-cli-file-organizer
go build -o fm .

On Windows, use go build -o fm.exe . instead. Requires Go 1.22 or newer.

Usage
fm organize [-apply] [dir]   sort files into folders by type (dry run without -apply)
fm undo [dir]                reverse organize
fm stats [dir]               count files and bytes per type
fm find pattern [dir]        list files matching a glob such as "*.pdf"
fm dupes [dir]               list files with identical content

If no folder is given, the current folder is used.

Organize a folder
bash
fm organize ~/Downloads          # show the plan, change nothing
fm organize -apply ~/Downloads   # move the files
fm undo ~/Downloads              # put everything back

Files go into Images, PDFs, Documents, Spreadsheets, Audio, Video, Archives, Installers, or Other. Only files at the top level of the folder are moved; subfolders are left alone.

See what's using space
bash
fm stats ~/Downloads
TYPE          FILES  SIZE
Video         1      7.6 MB
Audio         1      3.3 MB
Images        4      2.9 MB
PDFs          2      361.3 KB
Find files
bash
fm find "*.pdf" ~/Downloads

Matching ignores case, so *.jpg also finds IMG_01.JPG. Keep the quotes, or your shell expands the pattern first.

Find duplicates
bash
fm dupes ~/Pictures
3 copies, 878.9 KB each:
  IMG_2041.jpg
  vacation (1).jpg
  vacation.jpg

Duplicates are found by content, so files match even when their names are different.

Safety
organize is a dry run unless you pass -apply.
Existing files are never overwritten. A file whose destination is taken is left where it is.
Every move is recorded in .fm-journal.json inside the organized folder. Don't delete it before running undo; undo removes it when it's done.
If undo can't restore a file because something new is in its place, it stops and keeps the remaining moves, so running it again picks up where it left off.
Hidden files, and Windows system files like desktop.ini, are never touched.
How it works
organize plans every move before making any, so the dry run shows exactly what -apply will do. Moving a file is a rename, so even large files move instantly.
dupes groups files by size first, since files of different sizes can't match, and only hashes (SHA-256) the files that share a size. Most files are never read.
stats, find, and dupes share one directory walker that skips hidden folders like .git and keeps going past folders it can't read.
Development
bash
go vet ./...
go test ./...