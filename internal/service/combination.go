package service

import (
	"fmt"
	"github.com/vcsfrl/random-fit/internal/combination"
	"os"
	"path/filepath"
	"strings"
)

var ErrCombinationDefinitionManager = "combination definition manager error"

var DefinitionTemplate string

type CombinationStarDefinitionManager struct {
	dataFolder string
}

func NewCombinationStarDefinitionManager(dataFolder string) *CombinationStarDefinitionManager {
	return &CombinationStarDefinitionManager{
		dataFolder: dataFolder,
	}
}

func (dm *CombinationStarDefinitionManager) List() ([]string, error) {
	result := make([]string, 0)

	// print all files from the definitions folder
	files, err := os.ReadDir(dm.dataFolder)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", ErrCombinationDefinitionManager, err)
	}

	for _, file := range files {
		if file.IsDir() || file.Name()[0] == '.' {
			continue
		}

		result = append(result, strings.TrimSuffix(filepath.Base(file.Name()), filepath.Ext(file.Name())))
	}

	return result, nil
}

func (dm *CombinationStarDefinitionManager) New(definitionName string) error {
	// create a file for the definition
	definitionFileName := fmt.Sprintf("%s.star", definitionName)
	definitionFilePath := filepath.Join(dm.dataFolder, definitionFileName)

	// check if the file already exists
	if _, err := os.Stat(definitionFilePath); !os.IsNotExist(err) {
		return fmt.Errorf("%s: definition already exists", ErrCombinationDefinitionManager)
	}

	if err := os.WriteFile(definitionFilePath, []byte(DefinitionTemplate), 0644); err != nil {
		return fmt.Errorf("%s: new definition: %w", ErrCombinationDefinitionManager, err)
	}

	return nil
}

func (dm *CombinationStarDefinitionManager) GetScript(definitionName string) (string, error) {
	definitionFileName := fmt.Sprintf("%s.star", definitionName)
	definitionFilePath := filepath.Join(dm.dataFolder, definitionFileName)

	if _, err := os.Stat(definitionFilePath); os.IsNotExist(err) {
		return "", fmt.Errorf("%s: definition does not exist", ErrCombinationDefinitionManager)
	}

	return definitionFilePath, nil
}

func (dm *CombinationStarDefinitionManager) Build(definitionName string) (*combination.Combination, error) {
	definitionScript, err := dm.GetScript(definitionName)
	if err != nil {
		return nil, fmt.Errorf("%s: getting script: %w", ErrCombinationDefinitionManager, err)
	}

	definition, err := combination.NewCombinationDefinition(definitionScript)
	if err != nil {
		return nil, fmt.Errorf("%s: creating combination definition: %w", ErrCombinationDefinitionManager, err)
	}

	starBuilder, err := combination.NewStarBuilder(definition)
	if err != nil {
		return nil, fmt.Errorf("%s: creating star builder: %w", ErrCombinationDefinitionManager, err)
	}

	builtCombination, err := starBuilder.Build()
	if err != nil {
		return nil, fmt.Errorf("%s: building combination: %w", ErrCombinationDefinitionManager, err)
	}

	return builtCombination, nil
}

func (dm *CombinationStarDefinitionManager) Delete(name string) error {
	scriptName, err := dm.GetScript(name)
	if err != nil {
		return fmt.Errorf("%s: getting script: %w", ErrCombinationDefinitionManager, err)
	}

	if err := os.Remove(scriptName); err != nil {
		return fmt.Errorf("%s: deleting script: %w", ErrCombinationDefinitionManager, err)
	}

	return nil
}
