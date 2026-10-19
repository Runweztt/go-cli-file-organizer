package main

import (
	"fmt"
	"io/fs"
	"os"
	"sort"
	"text/tabwriter"
)

// runStats prints how many files and bytes each category holds under dir,
// largest category first.
func runStats(args []string) error {
	type total struct {
		name  string
		files int
		bytes int64
	}
	totals := map[string]*total{}

	err := walkFiles(dirArg(args, 0), func(_ string, info fs.FileInfo) error {
		name := category(info.Name())
		t := totals[name]
		if t == nil {
			t = &total{name: name}
			totals[name] = t
		}
		t.files++
		t.bytes += info.Size()
		return nil
	})
	if err != nil {
		return err
	}

	rows := make([]*total, 0, len(totals))
	for _, t := range totals {
		rows = append(rows, t)
	}
	sort.Slice(rows, func(i, j int) bool { return rows[i].bytes > rows[j].bytes })

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "TYPE\tFILES\tSIZE")
	for _, t := range rows {
		fmt.Fprintf(w, "%s\t%d\t%s\n", t.name, t.files, humanSize(t.bytes))
	}
	return w.Flush()
}
