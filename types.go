package main

import (
	"io"
	"net/http"
	"os"
)

type ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request) error

type ProblemYaml struct {
	TestsPath string `yaml:"tests_path"`
}

func readFileBytes(pathName string) ([]byte, error) {
	file, err := os.Open(pathName)

	if err != nil {
		return nil, err
	}

	return io.ReadAll(file)
}
