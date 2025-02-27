package main

//import "os/exec"

type Tester interface {
	Run(string) ([]int, error)
}

type CppTester struct {
	Problem Problem
}

func NewCppTester(problem Problem) *CppTester {
	return &CppTester{Problem: problem}
}

func (c *CppTester) Run(binPathName string) ([]int, error) {

	test := c.Problem.NextTestCase()
	results := []int{}
	for test != nil {
		result, err := test.RunTestCase(binPathName)
		if err != nil {
			return results, err
		}
		results = append(results, result)
		test = c.Problem.NextTestCase()
	}
	//	cmd := exec.Command(c.binPathName)

	return results, nil
}
