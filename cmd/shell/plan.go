package shell

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var ErrPlanManager = "plan manager error"

type PlanManager struct {
	dataFolder string
}

func (m *PlanManager) List() ([]string, error) {
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

func NewPlanManager(folder string) *PlanManager {
	return &PlanManager{
		dataFolder: folder,
	}
}
