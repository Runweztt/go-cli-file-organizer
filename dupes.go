package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"io/fs"
	"os"
)

// runDupes prints groups of files under dir that have identical content.
//
// Only files that share a size can be identical, and sizes come free with
// the directory walk, so files are grouped by size first and only those
// groups are hashed. Hashing means reading every byte, so this skips most
// of the expensive work.
func runDupes(args []string) error {
	bySize := map[int64][]string{}
	err := walkFiles(dirArg(args, 0), func(path string, info fs.FileInfo) error {
		if info.Size() > 0 {
			bySize[info.Size()] = append(bySize[info.Size()], path)
		}
		return nil
	})
	if err != nil {
		return err
	}

	for size, paths := range bySize {
		if len(paths) < 2 {
			continue
		}
		byHash := map[[sha256.Size]byte][]string{}
		for _, p := range paths {
			sum, err := hashFile(p)
			if err != nil {
				fmt.Fprintln(os.Stderr, "fm: skipping:", err)
				continue
			}
			byHash[sum] = append(byHash[sum], p)
		}
		for _, same := range byHash {
			if len(same) < 2 {
				continue
			}
			fmt.Printf("%d copies, %s each:\n", len(same), humanSize(size))
			for _, p := range same {
				fmt.Println("  " + p)
			}
		}
	}
	return nil
}

func hashFile(path string) (sum [sha256.Size]byte, err error) {
	f, err := os.Open(path)
	if err != nil {
		return sum, err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return sum, err
	}
	copy(sum[:], h.Sum(nil))
	return sum, nil
}
