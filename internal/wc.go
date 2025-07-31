package internal

import (
	"bufio"
	"cmp"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"
)

type record struct {
	date string
	name string
	num  int
}

func newRecord(date string, name string, num int) record {
	return record{date: date, name: name, num: num}
}

func CmdWC(args []string) error {
	if len(args) != 2 {
		fmt.Println("NOTE: wc needs a filepath")
		fmt.Println("exit")
		os.Exit(0)
	}
	path := args[1]
	var records []record
	records, err := executeWC(path)
	if err != nil {
		return err
	}
	total := 0
	for _, r := range records {
		total += r.num
	}
	output(records, total)
	return nil
}

func executeWC(fpath string) ([]record, error) {
	var records []record
	err := filepath.WalkDir(fpath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if isSecret(d.Name()) {
			return nil
		}
		if d.IsDir() {
			return nil
		}
		r, err := recordFromFile(path)
		if err != nil {
			return nil
		}
		records = append(records, r)
		return nil
	})
	if err != nil {
		return nil, err
	}
	sortByDate(records)
	return records, nil
}

func recordFromFile(path string) (record, error) {
	f, err := os.Open(path)
	if err != nil {
		return record{}, err
	}
	info, err := os.Lstat(path)
	if err != nil {
		return record{}, err
	}
	r := newRecord(info.ModTime().Format(time.DateOnly), info.Name(), countWords(f))
	return r, nil
}

func countWords(f *os.File) int {
	num := 0
	sc := bufio.NewScanner(f)
	sc.Split(bufio.ScanWords)
	for sc.Scan() {
		if !isWord(sc.Text()) {
			continue
		}
		num++
	}
	return num
}

func cleanWord(word string) string {
	signs := []string{",", ".", "-", "?", "ー", "!"}
	for _, s := range signs {
		if strings.Contains(word, s) {
			word = strings.ReplaceAll(word, s, "")
		}
	}
	return word
}

func isWord(input string) bool {
	return cleanWord(input) != ""
}

func isSecret(s string) bool {
	return strings.HasPrefix(s, ".")
}

func sortByDate(records []record) []record {
	slices.SortFunc(records, func(a, b record) int {
		return cmp.Compare(a.date, b.date)
	})
	return records
}

func output(records []record, total int) {
	fmt.Println()
	fmt.Printf(" %s : %d words\n", "total", total)
	fmt.Println(" ----------------------------------------")
	fmt.Printf(" %s          %s         %s\n", "DATE", "FILE", "WORD COUNT")
	fmt.Println(" ----------------------------------------")
	for _, r := range records {
		fmt.Printf(" %s    %s        %d words\n", r.date, r.name, r.num)
	}
	fmt.Println()
}
