package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/bytebase/pgparser/nodes"
	"github.com/bytebase/pgparser/parser"
)

type result struct {
	pgOut     string
	pgErr     string
	pgparser  string
	parseErr  error
	pgExitErr error
}

func main() {
	var helperPath string
	var filePath string
	var dirPath string

	flag.StringVar(&helperPath, "pg-helper", "tools/pg_parse_helper/pg_parse_helper", "path to pg_parse_helper binary")
	flag.StringVar(&filePath, "file", "", "single SQL file to compare")
	flag.StringVar(&dirPath, "dir", "", "directory of .sql files to compare")
	flag.Parse()

	if filePath == "" && dirPath == "" {
		filePath = "-"
	}

	if dirPath != "" {
		if err := runDir(helperPath, dirPath); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
		return
	}

	if err := runFile(helperPath, filePath); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func runDir(helperPath, dirPath string) error {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return fmt.Errorf("read dir: %w", err)
	}
	var files []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if strings.HasSuffix(entry.Name(), ".sql") {
			files = append(files, filepath.Join(dirPath, entry.Name()))
		}
	}
	sort.Strings(files)
	if len(files) == 0 {
		return fmt.Errorf("no .sql files in %s", dirPath)
	}
	for _, file := range files {
		if err := runFile(helperPath, file); err != nil {
			return err
		}
	}
	return nil
}

func runFile(helperPath, filePath string) error {
	sql, err := readSQL(filePath)
	if err != nil {
		return err
	}
	res := compareSQL(helperPath, sql)
	if res.pgExitErr != nil {
		return fmt.Errorf("pg helper error: %v: %s", res.pgExitErr, res.pgErr)
	}
	if res.parseErr != nil {
		return fmt.Errorf("pgparser error: %v", res.parseErr)
	}
	if res.pgOut == res.pgparser {
		if filePath == "-" {
			fmt.Println("OK")
		} else {
			fmt.Printf("OK %s\n", filePath)
		}
		return nil
	}

	label := "stdin"
	if filePath != "-" {
		label = filePath
	}
	diffIdx := firstDiffIndex(res.pgOut, res.pgparser)
	fmt.Printf("MISMATCH %s\n", label)
	fmt.Printf("PG length: %d\n", len(res.pgOut))
	fmt.Printf("pgparser length: %d\n", len(res.pgparser))
	if diffIdx >= 0 {
		fmt.Printf("first diff index: %d\n", diffIdx)
		fmt.Printf("PG context: %s\n", snippet(res.pgOut, diffIdx))
		fmt.Printf("pgparser context: %s\n", snippet(res.pgparser, diffIdx))
	}
	return fmt.Errorf("nodeToString mismatch")
}

func compareSQL(helperPath, sql string) result {
	pgOut, pgErr, pgExitErr := runPGHelper(helperPath, sql)
	pgparserOut, parseErr := runPGParser(sql)
	return result{
		pgOut:     pgOut,
		pgErr:     pgErr,
		pgparser:  pgparserOut,
		parseErr:  parseErr,
		pgExitErr: pgExitErr,
	}
}

func runPGHelper(helperPath, sql string) (string, string, error) {
	cmd := exec.Command(helperPath)
	cmd.Stdin = strings.NewReader(sql)
	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf
	err := cmd.Run()
	return strings.TrimSpace(outBuf.String()), strings.TrimSpace(errBuf.String()), err
}

func runPGParser(sql string) (string, error) {
	list, err := parser.Parse(sql)
	if err != nil {
		return "", err
	}
	return nodes.NodeToString(list), nil
}

func readSQL(path string) (string, error) {
	if path == "-" {
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			return "", fmt.Errorf("read stdin: %w", err)
		}
		return string(data), nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("read file %s: %w", path, err)
	}
	return string(data), nil
}

func firstDiffIndex(a, b string) int {
	max := len(a)
	if len(b) < max {
		max = len(b)
	}
	for i := 0; i < max; i++ {
		if a[i] != b[i] {
			return i
		}
	}
	if len(a) != len(b) {
		return max
	}
	return -1
}

func snippet(s string, idx int) string {
	if idx < 0 {
		return ""
	}
	start := idx - 40
	if start < 0 {
		start = 0
	}
	end := idx + 40
	if end > len(s) {
		end = len(s)
	}
	return s[start:end]
}
