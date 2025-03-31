package main

import (
	"bufio"
	"fmt"
	"github.com/subpop/go-log"
	"os"
	"strings"
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
		closeErr := file.Close()
		if closeErr != nil {
			log.Errorf("cannot close file %s: %s", name, closeErr)
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
func writeLines(name string, lines []string, truncate bool) (err error) {
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

		// When it is not possible to close the file, then return the error
		// reported by the failed close() system call.
		// This is very important in the situation when the system call close()
		// returns EIO. This can be caused when the previously-uncommitted write(2)
		// encountered an input/output error, resulting in the possible loss of data.
		defer func() {
			closeErr := file.Close()
			if closeErr != nil {
				err = closeErr
			}
		}()

		if _, err := file.Write(data); err != nil {
			return fmt.Errorf("cannot write to file: %w", err)
		}
	}
	return err
}
