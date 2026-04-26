package fs

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const FolderPermission = 0755
const FilePermission = 0644

func CreateFolder(folder string) error {
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		if err := os.MkdirAll(folder, FolderPermission); err != nil {
			return fmt.Errorf("error creating folder %s: %w", folder, err)
		}
	}

	return nil
}

// ListFileNames reads a directory and returns a list of file names without extensions,
// skipping directories and hidden files (starting with '.').
func ListFileNames(folder string) ([]string, error) {
	result := make([]string, 0)

	files, err := os.ReadDir(folder)
	if err != nil {
		return nil, fmt.Errorf("error reading folder %s: %w", folder, err)
	}

	for _, file := range files {
		if file.IsDir() || file.Name()[0] == '.' {
			continue
		}

		result = append(result, strings.TrimSuffix(filepath.Base(file.Name()), filepath.Ext(file.Name())))
	}

	return result, nil
}
