package shell

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var ErrDefinitionManager = "definition manager error"

type StarDefinitionManager struct {
	dataFolder string
}

func NewStarDefinitionManager(dataFolder string) *StarDefinitionManager {
	return &StarDefinitionManager{
		dataFolder: dataFolder,
	}
}

func (dm *StarDefinitionManager) List() ([]string, error) {
	result := make([]string, 0)

	// print all files from the definitions folder
	files, err := os.ReadDir(dm.dataFolder)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", ErrDefinitionManager, err)
	}

	for _, file := range files {
		if file.IsDir() || file.Name()[0] == '.' {
			continue
		}

		result = append(result, strings.TrimSuffix(filepath.Base(file.Name()), filepath.Ext(file.Name())))
	}

	return result, nil
}

func (dm *StarDefinitionManager) New(definitionName string) error {
	// create a file for the definition
	definitionFileName := fmt.Sprintf("%s.star", definitionName)
	definitionFilePath := filepath.Join(dm.dataFolder, definitionFileName)

	// check if the file already exists
	if _, err := os.Stat(definitionFilePath); !os.IsNotExist(err) {
		return fmt.Errorf("%s: definition already exists", ErrDefinitionManager)
	}

	if err := os.WriteFile(definitionFilePath, []byte(definitionTemplate), 0644); err != nil {
		return fmt.Errorf("%s: new definition: %w", ErrDefinitionManager, err)
	}

	return nil
}

func (dm *StarDefinitionManager) GetScript(definitionName string) (string, error) {
	definitionFileName := fmt.Sprintf("%s.star", definitionName)
	definitionFilePath := filepath.Join(dm.dataFolder, definitionFileName)

	if _, err := os.Stat(definitionFilePath); os.IsNotExist(err) {
		return "", fmt.Errorf("%s: definition does not exist", ErrDefinitionManager)
	}

	return definitionFilePath, nil
}
