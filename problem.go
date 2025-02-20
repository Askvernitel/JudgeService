package main

import (
	"bytes"
	_ "bytes"
	"fmt"
	"os"
	"path/filepath"

	yaml "gopkg.in/yaml.v3"
)

type Problem interface {
	NextTestCase() TestCase
}

//var problemsPath string = os.Getenv("PROBLEMS_PATH")

const (
	PROBLEM_YAML = "problem.yaml"
)

type NormalProblem struct {
	ProblemPathName  string
	TestsDirPathName string
	TestCases        []*ProblemTestCase
}

func NewNormalProblem(contestName, problemName string) *NormalProblem {

	return &NormalProblem{ProblemPathName: fmt.Sprintf("%s/%s/%s", problemsPath, contestName, problemName), TestCases: []*ProblemTestCase{}}
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
func (p *NormalProblem) addTestCases() error {
	fullTestDirPath := fmt.Sprintf("%s/%s", p.ProblemPathName, p.TestsDirPathName)
	dirs, err := os.ReadDir(fullTestDirPath)
	if err != nil {
		return err
	}
	fmt.Println(fullTestDirPath)
	for _, file := range dirs {
		if file.IsDir() {
			dirName := file.Name()
			fmt.Println(dirName)
			files, err := os.ReadDir(fmt.Sprintf("%s/%s", fullTestDirPath, dirName))
			if err != nil {
				return err
			}
			//TODO: Separate this into functions
			var inFilePath, outFilePath string
			for _, inOutFile := range files {
				fmt.Println("inOutFile: " + inOutFile.Name())
				if !inOutFile.IsDir() && filepath.Ext(inOutFile.Name()) == "in" && inFilePath == "" {
					inFilePath = fmt.Sprintf("%s/%s/%s", fullTestDirPath, dirName, inOutFile.Name())
					fmt.Println(inFilePath)
				}
				if !inOutFile.IsDir() && filepath.Ext(inOutFile.Name()) == "ans" && outFilePath == "" {
					outFilePath = fmt.Sprintf("%s/%s/%s", fullTestDirPath, dirName, inOutFile.Name())
					fmt.Println(inFilePath)
				}
			}
			if inFilePath == "" || outFilePath == "" {
				continue
			}
			p.TestCases = append(p.TestCases, NewProblemTestCase(inFilePath, outFilePath))
		}
	}
	return nil
}
func (p *NormalProblem) initProblemTestCases() error {
	//dirFiles, err := os.ReadDir(p.problemPathName)
	problemInfo, err := p.readProblemYaml()
	if err != nil {
		return err
	}
	p.TestsDirPathName = problemInfo.TestsPath
	err = p.addTestCases()
	if err != nil {
		return err
	}
	for _, test := range p.TestCases {
		fmt.Println(test.TestInput)
	}
	return nil
}
func (p *NormalProblem) NextTestCase() TestCase {
	return nil
}

type TestCase interface {
	RunTestCase(string) int
}

type ProblemTestCase struct {
	TestInput  string
	TestOutput string
}

func NewProblemTestCase(testInput string, testOutput string) *ProblemTestCase {
	return &ProblemTestCase{TestInput: testInput, TestOutput: testOutput}
}

func (t *ProblemTestCase) RunTestCase() {
}
