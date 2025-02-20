package main

import "net/http"

type ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request) error

type ProblemYaml struct {
	TestsPath string `yaml:"tests_path"`
}
