package main

//import "os/exec"

type Tester interface {
	Run(string) ([]*TestResult, error)
}

type CppTester struct {
	Problem Problem
}

func NewCppTester(problem Problem) *CppTester {
	return &CppTester{Problem: problem}
}

//problemLimits := c.Problem.GetTestLimits()
//good Limiters commented
//	limiter := NewCmdLimiter(binPathName, problemLimits.MemoryLimitMb, problemLimits.TimeLimitSec)

//limiter := NewNoLimiter(binPathName)

// cmd := exec.Command(c.binPathName)
func (c *CppTester) Run(binPathName string) ([]*TestResult, error) {
	problemLimits := c.Problem.GetTestLimits()
	limiter := NewCmdLimiter(binPathName, problemLimits.MemoryLimitMb, problemLimits.TimeLimitSec)
	test := c.Problem.NextTestCase()
	results := []*TestResult{}
	for test != nil {
		result, err := test.RunTestCase(binPathName, limiter)

		if err != nil {
			return results, err
		}
		results = append(results, result)
		test = c.Problem.NextTestCase()
	}

	return results, nil
}
