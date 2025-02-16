package main

//import "os/exec"

type Tester interface {
	Run() (int, error)
}

type CppTester struct {
	binPathName string
}

func NewCppTester(binPathName string, problem Problem) *CppTester {
	return &CppTester{binPathName: binPathName}
}

func (c *CppTester) Run() error {
	//	cmd := exec.Command(c.binPathName)

	return nil
}
