package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"git.sr.ht/~spc/go-log"
)

// canonicalizeRepoName converts the string filename into a suitable filename.
// It replaces all spaces with underscores and appends the given suffix.
func canonicalizeRepoName(filename, suffix string) string {
	return strings.TrimSuffix(strings.ReplaceAll(filename, " ", "_"), suffix) + suffix
}

// readLines reads each line in file into a slice.
func readLines(name string) ([]string, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, fmt.Errorf("cannot open file for reading: %w", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Warnf("cannot close file %v: %v", name, err)
		}
	}()

	var lines []string
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("cannot read file: %w", err)
	}

	return lines, nil
}

// writeLines writes lines to file name. If truncate is true, the file is
// truncated before writing.
func writeLines(name string, lines []string, truncate bool) error {
	var data []byte

	for _, line := range lines {
		data = append(data, line+"\n"...)
	}

	if truncate {
		return os.WriteFile(name, data, 0644)
	} else {
		file, err := os.OpenFile(name, os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return fmt.Errorf("cannot open file for appending: %w", err)
		}
		defer func() {
			if err := file.Close(); err != nil {
				log.Warnf("cannot close file %v: %v", name, err)
			}
		}()

		if _, err := file.Write(data); err != nil {
			return fmt.Errorf("cannot write to file: %w", err)
		}
	}
	return nil
}
