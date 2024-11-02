// Package main bootstrap the analyzer to validate the project structure.
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
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

func configureLogging(loggerLevel string) error {
	logrusLoggerLevel, err := log.ParseLevel(loggerLevel)
	if err != nil {
		return fmt.Errorf("unable to parse the logrusLoggerLevel: '%w'", err)
	}
	log.SetLevel(logrusLoggerLevel)
	log.SetReportCaller(true)
	return nil
}

func main() {
	loggerLevel := flag.String("loggerLevel", "Fatal", "set the loggerLevel to either: Trace, Debug, Info, Warning, Error, Fatal or Panic")
	flag.Parse()

	if err := configureLogging(*loggerLevel); err != nil {
		log.WithError(err).Fatal("could not configure logging")
	}

	unmarshalledProLayoutYML, err := readAndUnmarshalProLayoutYML(proLayoutFile)
	if err != nil {
		log.WithError(err).Fatal("failed to unmarshal")
	}

	singlechecker.Main(analyzer.New(*unmarshalledProLayoutYML))
}
