package main

import (
	"fmt"
	"os"
)

var problemsPath string = os.Getenv("PROBLEMS_PATH")

type Problem struct {
	testCases       []string
	problemPathName string
}

func NewProblem(contestName, problemName string) *Problem {

	return &Problem{problemPathName: fmt.Sprintf("%s/%s/%s", problemsPath, contestName, problemName)}
}

func (p *Problem) initProblemTestCases() {

}

type ProblemTestCase struct {
}
