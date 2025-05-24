package fs

import (
	"fmt"
	"os"
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
