// Package main bootstrap the analyzer to validate the project structure.
package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"golang.org/x/tools/go/analysis/singlechecker"
	"gopkg.in/yaml.v3"

	"github.com/wimspaargaren/prolayout/internal/analyzer"
	"github.com/wimspaargaren/prolayout/internal/model"
)

const proLayoutFile = ".prolayout.yml"

func readAndUnmarshalProLayoutYML(proLayoutFile string) (*model.Root, error) {
	data, err := os.ReadFile(filepath.Clean(proLayoutFile))
	if err != nil {
		return nil, fmt.Errorf("'%w'", err)
	}
	t := model.Root{}
	err = yaml.Unmarshal(data, &t)
	if err != nil {
		return nil, fmt.Errorf("'%w'", err)
	}

	return &t, nil
}

func main() {
	unmarshalledProLayoutYML, err := readAndUnmarshalProLayoutYML(proLayoutFile)
	if err != nil {
		log.Fatalf("failed to unmarshal '%s': '%v'", proLayoutFile, err)
	}

	singlechecker.Main(analyzer.New(*unmarshalledProLayoutYML))
}
