package plan

import (
	"encoding/gob"
	"fmt"
	"github.com/vcsfrl/random-fit/internal/combination"
	"os"
	"path/filepath"
	"strings"
)

var ErrExport = fmt.Errorf("error exporting plan")

type Exporter struct {
	OutputDir  string
	StorageDir string

	workers int
}

func NewExporter(outputDir string, storageDir string) *Exporter {
	return &Exporter{
		OutputDir:  outputDir,
		StorageDir: storageDir,
	}
}

func (e *Exporter) Export(plan *UserPlan) error {
	for userID, groups := range plan.UserGroups {
		for _, group := range groups {
			groupFolder := strings.ReplaceAll(filepath.Join(e.OutputDir, userID, e.containerFolder(plan.Plan, group.Group), group.Details), " ", "_")
			if err := os.MkdirAll(groupFolder, 0755); err != nil {
				return fmt.Errorf("%w: error creating group folder: %s", ErrExport, err)
			}

			// Create a file for each combination by type
			for i, groupCombination := range group.Combinations {
				for _, data := range groupCombination.Data {
					if err := e.saveToFile(groupCombination, data, groupFolder, i+1); err != nil {
						return err
					}
				}
			}
		}
	}

	if err := e.exportObject(plan); err != nil {
		return fmt.Errorf("%w: error exporting plan object: %s", ErrExport, err)
	}

	return nil
}

func (e *Exporter) ExportGenerator(generator chan *PlannedCombination) error {
	for planCombination := range generator {
		if planCombination.Err != nil {
			return fmt.Errorf("%w: error generating plan: %s", ErrExport, planCombination.Err)
		}

		groupFolder := strings.ReplaceAll(filepath.Join(e.OutputDir, planCombination.User, e.containerFolder(planCombination.Plan, planCombination.Group), planCombination.Group.Details), " ", "_")
		if err := os.MkdirAll(groupFolder, 0755); err != nil {
			return fmt.Errorf("%w: error creating group folder: %s", ErrExport, err)
		}

		// Create a file for each combination by type
		for _, data := range planCombination.Combination.Data {
			if err := e.saveToFile(planCombination.Combination, data, groupFolder, planCombination.GroupSerialId); err != nil {
				return fmt.Errorf("%w: error saving file: %s", ErrExport, err)
			}
		}

		if err := e.exportPlannedCombinationObject(planCombination); err != nil {
			return fmt.Errorf("%w: error exporting plan object: %s", ErrExport, err)
		}
	}

	return nil
}

func (e *Exporter) containerFolder(plan Plan, group Group) string {
	if len(group.ContainerName) == 0 {
		return plan.DefinitionID
	}

	var folder string

	for _, container := range group.ContainerName {
		if container == "_date" {
			folder = filepath.Join(folder, plan.CreatedAt.Format("2006-01-02-15-04"))
			continue
		}
		folder = filepath.Join(folder, container)
	}

	return folder
}

func (e *Exporter) exportObject(plan *UserPlan) error {
	// save the plan to storage
	storageFile := filepath.Join(e.StorageDir, fmt.Sprintf("%s.gob", plan.UUID.String()))
	//open the file
	file, err := os.Create(storageFile)
	if err != nil {
		return fmt.Errorf("%w: error creating storage file: %s", ErrExport, err)
	}

	defer func() {
		_ = file.Close()
	}()

	encoder := gob.NewEncoder(file)
	if err := encoder.Encode(plan); err != nil {
		return fmt.Errorf("%w: error encoding plan object: %s", ErrExport, err)
	}

	return nil
}

func (e *Exporter) exportPlannedCombinationObject(plan *PlannedCombination) error {
	// save the plan to storage
	storageFile := filepath.Join(e.StorageDir, fmt.Sprintf("%s_%s_%s.gob", plan.User, plan.UUID.String(), plan.Combination.UUID.String()))
	//open the file
	file, err := os.Create(storageFile)
	if err != nil {
		return fmt.Errorf("%w: error creating storage file: %s", ErrExport, err)
	}

	defer func() {
		_ = file.Close()
	}()

	encoder := gob.NewEncoder(file)
	if err := encoder.Encode(plan); err != nil {
		return fmt.Errorf("%w: error encoding plan object: %s", ErrExport, err)
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
