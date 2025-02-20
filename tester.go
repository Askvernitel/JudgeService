package main

//import "os/exec"

type Tester interface {
	Run(string) (int, error)
}

type CppTester struct {
}

func NewCppTester(problem Problem) *CppTester {
	return &CppTester{}
}

func (c *CppTester) Run(binPathName string) (int, error) {
	//	cmd := exec.Command(c.binPathName)

	return 0, nil
}
