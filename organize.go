package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// journalName is the file, kept inside the organized directory, that records
// every move so undo can reverse it. It is hidden, so organize ignores it.
const journalName = ".fm-journal.json"

// A move is a rename, with both paths relative to the organized directory.
// Relative paths keep the journal valid if the directory itself is moved.
type move struct {
	From string `json:"from"`
	To   string `json:"to"`
}

// runOrganize moves each top-level file of dir into a folder named after its
// category. Without -apply it only prints the plan.
func runOrganize(args []string) error {
	flags := flag.NewFlagSet("organize", flag.ExitOnError)
	apply := flags.Bool("apply", false, "move the files instead of printing the plan")
	flags.Parse(args)
	dir := dirArg(flags.Args(), 0)

	moves, err := plan(dir)
	if err != nil {
		return err
	}
	if len(moves) == 0 {
		fmt.Println("nothing to organize")
		return nil
	}
	if !*apply {
		for _, m := range moves {
			fmt.Printf("%s -> %s\n", m.From, m.To)
		}
		fmt.Printf("%d files to move; run again with -apply to move them\n", len(moves))
		return nil
	}

	journal, err := readJournal(dir)
	if err != nil {
		return err
	}
	moved := 0
	for _, m := range moves {
		to := filepath.Join(dir, m.To)
		if err = os.MkdirAll(filepath.Dir(to), 0o755); err != nil {
			break
		}
		if err = os.Rename(filepath.Join(dir, m.From), to); err != nil {
			break
		}
		journal = append(journal, m)
		moved++
	}

	// Record whatever moved, even after a failure, so it can still be undone.
	if moved > 0 {
		err = errors.Join(err, writeJournal(dir, journal))
	}
	fmt.Printf("moved %d files; run fm undo to reverse\n", moved)
	return err
}

// plan decides where each top-level file of dir should go without touching
// anything. A file whose destination is already taken is left where it is:
// os.Rename would silently overwrite it, and organize must never lose data.
func plan(dir string) ([]move, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var moves []move
	for _, e := range entries {
		if !e.Type().IsRegular() || hidden(e.Name()) {
			continue
		}
		to := filepath.Join(category(e.Name()), e.Name())
		if exists(filepath.Join(dir, to)) {
			fmt.Fprintf(os.Stderr, "fm: skipping %s: %s already exists\n", e.Name(), to)
			continue
		}
		moves = append(moves, move{From: e.Name(), To: to})
	}
	return moves, nil
}

// runUndo reverses the journaled moves in dir, newest first. On failure the
// journal keeps only the moves not yet undone, so undo can simply be rerun.
func runUndo(args []string) error {
	dir := dirArg(args, 0)
	journal, err := readJournal(dir)
	if err != nil {
		return err
	}
	if len(journal) == 0 {
		fmt.Println("nothing to undo")
		return nil
	}

	restored := 0
	for i := len(journal) - 1; i >= 0; i-- {
		from, to := filepath.Join(dir, journal[i].From), filepath.Join(dir, journal[i].To)
		switch {
		case !exists(to):
			fmt.Fprintf(os.Stderr, "fm: skipping %s: no longer there\n", to)
		case exists(from):
			err = fmt.Errorf("undo: %s already exists; move it away and run undo again", from)
		default:
			if err = os.Rename(to, from); err == nil {
				restored++
			}
		}
		if err != nil {
			return errors.Join(err, writeJournal(dir, journal[:i+1]))
		}
		os.Remove(filepath.Dir(to)) // only succeeds once the folder is empty
	}

	fmt.Printf("restored %d files\n", restored)
	return os.Remove(filepath.Join(dir, journalName))
}

func readJournal(dir string) ([]move, error) {
	data, err := os.ReadFile(filepath.Join(dir, journalName))
	if errors.Is(err, fs.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var moves []move
	err = json.Unmarshal(data, &moves)
	return moves, err
}

func writeJournal(dir string, moves []move) error {
	data, err := json.MarshalIndent(moves, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, journalName), append(data, '\n'), 0o644)
}

// exists reports whether anything is at path. It uses Lstat so a symlink,
// even a broken one, counts: renaming onto it would destroy it.
func exists(path string) bool {
	_, err := os.Lstat(path)
	return err == nil
}
