package parsers_test

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"issuefinder/infra/parsers"
)

func TestGetParser(t *testing.T) {
	doc := "<NessusClientData_v2></NessusClientData_v2>"
	p := parsers.GetParser(doc)
	assert.Equal(t, "Nessus", p.GetName())
}

// Helper functions for the other tests
func getDoc(file string) string {
	currentWorkingDirectory, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	path := getProjectDir(currentWorkingDirectory)
	f, err := os.Open(filepath.Join(path, file))
	if err != nil {
		return ""
	}
	defer f.Close()
	result, _ := ioutil.ReadAll(f)
	return string(result)
}

func getProjectDir(path string) string {
	if filepath.Base(path) == "go-IssueFinder" {
		return filepath.Join(path, "testdata")
	}
	return getProjectDir(filepath.Dir(path))
}
