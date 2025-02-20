package main

//import "os/exec"

type Tester interface {
	Run(string) (int, error)
}
type TestCase interface {
	Run(string) int
}

type CppTester struct {
}

func NewCppTester(problem *Problem) *CppTester {
	return &CppTester{}
}

func (c *CppTester) Run(binPathName string) (int, error) {
	//	cmd := exec.Command(c.binPathName)

	return 0, nil
}
