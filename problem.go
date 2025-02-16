package main

type Problem struct {
	problemPathName string
}

func NewProblem(problemPathName string) *Problem {
	return &Problem{problemPathName: problemPathName}
}

func (p *Problem) initProblemTests() {

}
