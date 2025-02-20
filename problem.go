package main

import (
	"bytes"
	_ "bytes"
	"fmt"
	"os"

	yaml "gopkg.in/yaml.v3"
)

type Problem interface {
	NextTestCase()
}

//var problemsPath string = os.Getenv("PROBLEMS_PATH")

const (
	PROBLEM_YAML = "problem.yaml"
)

type NormalProblem struct {
	ProblemPathName  string
	TestCasePathName string
}

func NewNormalProblem(contestName, problemName string) *NormalProblem {

	return &NormalProblem{ProblemPathName: fmt.Sprintf("%s/%s/%s", problemsPath, contestName, problemName)}
}
func (p *NormalProblem) readProblemYaml() (*ProblemYaml, error) {
	rawBytes, err := os.ReadFile(fmt.Sprintf("%s/%s", p.ProblemPathName, PROBLEM_YAML))
	if err != nil {
		return nil, err
	}
	bytesReader := bytes.NewReader(rawBytes)
	decoder := yaml.NewDecoder(bytesReader)

	problemInfo := &ProblemYaml{}
	err = decoder.Decode(problemInfo)
	if err != nil {
		return nil, err
	}
	return problemInfo, nil

}
func (p *NormalProblem) initProblemTestCases() error {
	//dirFiles, err := os.ReadDir(p.problemPathName)
	problemInfo, err := p.readProblemYaml()
	if err != nil {
		return err
	}
	p.TestCasePathName = problemInfo.TestsPath
	return nil
}

type ProblemTestCase struct {
}
