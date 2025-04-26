package plan

import (
	"fmt"
	"github.com/vcsfrl/random-fit/internal/combination"
	"os"
	"path/filepath"
	"strings"
)

var ErrExport = fmt.Errorf("error exporting plan")

type Exporter struct {
	OutputDir string
}

func NewExporter(outputDir string) *Exporter {
	return &Exporter{
		OutputDir: outputDir,
	}
}

func (e *Exporter) Export(plan *Plan) error {
	for userID, groups := range plan.UserGroups {
		for _, group := range groups {
			groupFolder, err := e.createGroupFolder(plan, group, userID)
			if err != nil {
				return err
			}

			// Create a file for each combination by type
			for _, combination := range group.Combinations {
				for _, data := range combination.Data {
					err2 := e.saveToFile(combination, data, groupFolder)
					if err2 != nil {
						return err2
					}
				}

			}

		}
	}
	return nil
}

func (e *Exporter) saveToFile(combination *combination.Combination, data *combination.Data, groupFolder string) error {
	fileName := fmt.Sprintf("%s.%s", combination.Details, data.Extension)
	filePath := filepath.Join(groupFolder, strings.ReplaceAll(fileName, " ", "_"))
	err := os.WriteFile(filePath, data.Data.Bytes(), 0666)
	if err != nil {
		return fmt.Errorf("%w: error writing file: %s", ErrExport, err)
	}

	return nil
}

func (e *Exporter) createGroupFolder(plan *Plan, group *Group, userID string) (string, error) {
	// Create a folder for the group
	groupFolder := filepath.Join(
		e.OutputDir,
		userID,
		fmt.Sprintf("%s_%s", plan.CreatedAt.Format("2006-01-02-1504"), group.ContainerName),
		group.Details,
	)

	groupFolder = strings.ReplaceAll(groupFolder, " ", "_")
	err := os.MkdirAll(groupFolder, 0755)
	if err != nil {
		return "", fmt.Errorf("%w: error creating group folder: %s", ErrExport, err)
	}

	return groupFolder, err
}
