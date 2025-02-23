package main

//import "os/exec"

type Tester interface {
	Run(string) (int, error)
}

type CppTester struct {
	Problem Problem
}

func NewCppTester(problem Problem) *CppTester {
	return &CppTester{Problem: problem}
}

func (c *CppTester) Run(binPathName string) (int, error) {

	test := c.Problem.NextTestCase()
	for test != nil {
		result, err := test.RunTestCase(binPathName)
		if err != nil {
			return result, nil
		}
		test = c.Problem.NextTestCase()
	}
	//	cmd := exec.Command(c.binPathName)

	return 0, nil
}
