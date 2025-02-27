package main

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"net/http"
	"os"
)

type ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request) error

type ProblemYaml struct {
	TestsPath string `yaml:"tests_path"`
}
type JudgeResponse struct {
	Results []int `json:"results"`
}

func readFileBytesByPath(pathName string) ([]byte, error) {
	file, err := os.Open(pathName)

	if err != nil {
		return nil, err
	}

	return io.ReadAll(file)
}
func CompareReaders(reader1, reader2 io.Reader) bool {
	scanner1 := bufio.NewScanner(reader1)
	scanner2 := bufio.NewScanner(reader2)
	for scanner1.Scan() && scanner2.Scan() {
		log.Println(scanner1.Text())
		if !bytes.Equal(scanner1.Bytes(), scanner2.Bytes()) {
			return false
		}
	}

	return !(scanner1.Scan() && scanner2.Scan())
}
