package plan

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

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
				return err
			}

			// Create a file for each combination by type
			for _, combination := range group.Combinations {
				for _, data := range combination.Data {
					// create file
					fileName := fmt.Sprintf("%s.%s", combination.Details, data.Extension)
					fileName = strings.ReplaceAll(fileName, " ", "_")
					filePath := filepath.Join(groupFolder, fileName)
					file, err := os.Create(filePath)
					if err != nil {
						return err
					}

					_, err = file.Write(data.Data.Bytes())
					if err != nil {
						return err

					}

					err = file.Close()
				}

			}

		}
	}
	return nil
}
