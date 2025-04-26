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
			groupFolder := strings.ReplaceAll(filepath.Join(e.OutputDir, userID, fmt.Sprintf("%s_%s", plan.CreatedAt.Format("2006-01-02-1504"), group.ContainerName), group.Details), " ", "_")
			if err := os.MkdirAll(groupFolder, 0755); err != nil {
				return fmt.Errorf("%w: error creating group folder: %s", ErrExport, err)
			}

			// Create a file for each combination by type
			for i, groupCombination := range group.Combinations {
				for _, data := range groupCombination.Data {
					if err := e.saveToFile(groupCombination, data, groupFolder, i); err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}

func (e *Exporter) saveToFile(groupCombination *combination.Combination, data *combination.Data, groupFolder string, i int) error {
	fileName := fmt.Sprintf("%s_%d.%s", groupCombination.Details, i, data.Extension)
	filePath := filepath.Join(groupFolder, strings.ReplaceAll(fileName, " ", "_"))
	err := os.WriteFile(filePath, data.Data.Bytes(), 0666)
	if err != nil {
		return fmt.Errorf("%w: error writing file: %s", ErrExport, err)
	}

	return nil
}
