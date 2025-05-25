package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	rfPlan "github.com/vcsfrl/random-fit/internal/plan"
	"github.com/vcsfrl/random-fit/internal/platform/fs"
	"os"
	"path/filepath"
	"strings"
)

var ErrPlanDefinitionManager = errors.New("plan definition manager error")

type PlanDefinitionManager struct {
	dataFolder string
}

func NewPlanDefinitionManager(folder string) *PlanDefinitionManager {
	return &PlanDefinitionManager{
		dataFolder: folder,
	}
}

func (m *PlanDefinitionManager) List() ([]string, error) {
	result := make([]string, 0)

	files, err := os.ReadDir(m.dataFolder)
	if err != nil {
		return nil, fmt.Errorf("%w: read data folder: %w", ErrPlanDefinitionManager, err)
	}

	for _, file := range files {
		if file.IsDir() || file.Name()[0] == '.' {
			continue
		}

		result = append(result, strings.TrimSuffix(filepath.Base(file.Name()), filepath.Ext(file.Name())))
	}

	return result, nil
}

func (m *PlanDefinitionManager) New(plan string) error {
	planFileName := plan + ".json"
	planFilePath := filepath.Join(m.dataFolder, planFileName)

	if _, err := os.Stat(planFilePath); !os.IsNotExist(err) {
		return fmt.Errorf("%w: plan already exists", ErrPlanDefinitionManager)
	}

	emptyPlan := m.GetSamplePlanDefinition()

	buff, err := json.Marshal(emptyPlan)
	if err != nil {
		return fmt.Errorf("%w: marshal plan to json: %w", ErrPlanDefinitionManager, err)
	}

	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, buff, "", "  "); err != nil {
		return fmt.Errorf("%w: indent json: %w", ErrPlanDefinitionManager, err)
	}

	if err := os.WriteFile(planFilePath, prettyJSON.Bytes(), fs.FilePermission); err != nil {
		return fmt.Errorf("%w: new plan: %w", ErrPlanDefinitionManager, err)
	}

	return nil
}

func (m *PlanDefinitionManager) GetFile(plan string) (string, error) {
	planFileName := plan + ".json"
	planFilePath := filepath.Join(m.dataFolder, planFileName)

	if _, err := os.Stat(planFilePath); os.IsNotExist(err) {
		return "", fmt.Errorf("%w: plan does not exist", ErrPlanDefinitionManager)
	}

	return planFilePath, nil
}

func (m *PlanDefinitionManager) Delete(name string) error {
	scriptName, err := m.GetFile(name)
	if err != nil {
		return fmt.Errorf("%w: get script: %w", ErrPlanDefinitionManager, err)
	}

	if err := os.Remove(scriptName); err != nil {
		return fmt.Errorf("%w: remove script: %w", ErrPlanDefinitionManager, err)
	}

	return nil
}

func (m *PlanDefinitionManager) GetSamplePlanDefinition() *rfPlan.Definition {
	return &rfPlan.Definition{
		ID:      "definition",
		Details: "Definition",
		Users:   []string{"user1"},
		UserData: rfPlan.UserData{
			ContainerName:            []string{"ContainerName", "_date"},
			RecurrentGroupNamePrefix: "Group",
			RecurrentGroups:          1,
			NrOfGroupCombinations:    1,
		},
	}
}
