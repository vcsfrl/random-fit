package shell

import (
	"bytes"
	"encoding/json"
	"fmt"
	rfPlan "github.com/vcsfrl/random-fit/internal/plan"
	"os"
	"path/filepath"
	"strings"
)

var ErrPlanManager = "plan manager error"

type PlanDefinitionManager struct {
	dataFolder string
}

func (m *PlanDefinitionManager) List() ([]string, error) {
	result := make([]string, 0)
	files, err := os.ReadDir(m.dataFolder)
	if err != nil {
		return nil, fmt.Errorf("%s: read data folder: %w", ErrPlanManager, err)
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
	planFileName := fmt.Sprintf("%s.json", plan)
	planFilePath := filepath.Join(m.dataFolder, planFileName)

	if _, err := os.Stat(planFilePath); !os.IsNotExist(err) {
		return fmt.Errorf("%s: plan already exists", ErrPlanManager)
	}

	emptyPlan := m.getSamplePlanDefinition()
	buff, err := json.Marshal(emptyPlan)
	if err != nil {
		return fmt.Errorf("%s: marshal plan to json: %w", ErrPlanManager, err)
	}

	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, buff, "", "  "); err != nil {
		return fmt.Errorf("%s: indent json: %w", ErrPlanManager, err)
	}

	if err := os.WriteFile(planFilePath, prettyJSON.Bytes(), 0644); err != nil {
		return fmt.Errorf("%s: new plan: %w", ErrPlanManager, err)
	}

	return nil
}

func (m *PlanDefinitionManager) getSamplePlanDefinition() *rfPlan.Definition {
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

func (m *PlanDefinitionManager) GetFile(plan string) (string, error) {
	planFileName := fmt.Sprintf("%s.json", plan)
	planFilePath := filepath.Join(m.dataFolder, planFileName)

	if _, err := os.Stat(planFilePath); os.IsNotExist(err) {
		return "", fmt.Errorf("%s: plan does not exist", ErrPlanManager)
	}

	return planFilePath, nil
}

func NewPlanDefinitionManager(folder string) *PlanDefinitionManager {
	return &PlanDefinitionManager{
		dataFolder: folder,
	}
}
