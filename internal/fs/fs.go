package fs

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/eslam-allam/file-signer/internal/slice"
)

type pathtype int

const (
	File pathtype = iota
	Directory
)

func Exists(path string) (bool, pathtype, error) {
	info, err := os.Stat(path)
	if err == nil {
		var t pathtype
		if info.IsDir() {
			t = Directory
		} else {
			t = File
		}
		return true, t, nil
	}
	if errors.Is(err, fs.ErrNotExist) {
		return false, 0, nil
	}

	return false, 0, err
}

func SaveCreateIntermediate(path string, bytes []byte, overwrite bool) error {
	exists, typ, err := Exists(path)
	if err != nil {
		return err
	}
	if exists {
		if typ != File {
			return fmt.Errorf("path '%s' already exists and is a directory", path)
		}
		if !overwrite {
			return fmt.Errorf("path '%s' already exists and overwriting is not permitted", path)
		}
	}
	directory := filepath.Dir(path)
	err = os.MkdirAll(directory, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to make intermediate directories for path '%s': %w", path, err)
	}
	err = os.WriteFile(path, bytes, 0644)
	if err != nil {
		return fmt.Errorf("failed to write file '%s': %w", path, err)
	}
	return nil
}

func ReadFile(path string) ([]byte, error) {
	exists, typ, err := Exists(path)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, fmt.Errorf("specified file '%s' does not exist", path)
	}
	if typ != File {
		return nil, fmt.Errorf("specified path '%s' is not a file", path)
	}
	return os.ReadFile(path)
}

func getFilesFilter(path, startsWith string, endsWith []string) []string {
	files, err := os.ReadDir(path)
	if err != nil {
		return []string{}
	}

	entries := make([]string, 0)
	// Iterate over the files and print their names
	for _, file := range files {
		fullPath := filepath.Join(path, file.Name())
		if file.IsDir() {
			entries = append(endsWith, getFilesFilter(fullPath, startsWith, endsWith)...)
			continue
		}
		if len(endsWith) != 0 && !slice.AnyMatch(endsWith,
			func(s string) bool {
				return strings.HasSuffix(file.Name(), s) || !strings.HasPrefix(file.Name(), startsWith)
			}) {
			continue
		}
		entries = append(entries, fullPath)
	}
	return entries
}

func ListDirFilter(path, startsWith string, endsWith []string) ([]string, error) {
	exists, typ, err := Exists(path)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, fmt.Errorf("directory '%s' does not exist", path)
	}
	if typ != Directory {
		return nil, fmt.Errorf("path '%s' is not a directory", path)
	}

	return getFilesFilter(path, startsWith, endsWith), nil
}
