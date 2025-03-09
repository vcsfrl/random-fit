package config

import (
	"fmt"
	"github.com/vcsfrl/random-fit/internal/core"
	"log"
	"os"
	"path/filepath"
)

func ElementsFromFolder(folderName string) []core.ElementDefinition {

	files, err := os.ReadDir(folderName)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		configFile := filepath.Join(folderName, f.Name())
		fmt.Println(configFile)
	}

	return []core.ElementDefinition{}
}
