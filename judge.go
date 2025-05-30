package main

import "fmt"

type Judge interface {
	Run()
}

const (
	RESULT_ACCEPTED              = 1
	RESULT_WRONG_ANSWER          = 2
	RESULT_TIME_EXCEEDED_LIMIT   = 3
	RESULT_COMPILATION_ERROR     = 4
	RESULT_RUNTIME_ERROR         = 5
	RESULT_JUDGE_ERROR           = 6
	RESULT_MEMORY_EXCEEDED_LIMIT = 7
)

// c++ judge
type CppJudge struct {
	Tester   Tester
	Compiler Compiler
	Results  []*TestResult
}

func NewCppJudge(compiler Compiler, tester Tester) *CppJudge {
	return &CppJudge{Compiler: compiler, Tester: tester, Results: []*TestResult{}}
}

func (j *CppJudge) Run() error {
	err := j.Compiler.Compile()
	fmt.Println("HERE: ", err)
	if err != nil {
		return err
	}
	binPathName := j.Compiler.OutputBinPathName()
	j.Results, err = j.Tester.Run(binPathName)
	fmt.Println(j.Results)
	return err
}

// java judge
type JavaJudge struct {
}

type PythonJudge struct {
}
