package main

import (
	"bytes"
	_ "bytes"
	"fmt"
	_ "io"
	"os"
	"os/exec"
	"path/filepath"

	yaml "gopkg.in/yaml.v3"
)

type Problem interface {
	NextTestCase() TestCase
	GetAllTestCases() []*ProblemTestCase
}

//var problemsPath string = os.Getenv("PROBLEMS_PATH")

const (
	PROBLEM_YAML = "problem.yaml"
)

type NormalProblem struct {
	ProblemPathName  string
	TestsDirPathName string
	TestCases        []*ProblemTestCase
	currentTestIndex int
}

func NewNormalProblem(contestName, problemName string) *NormalProblem {

	return &NormalProblem{ProblemPathName: fmt.Sprintf("%s/%s/%s", problemsPath, contestName, problemName), TestCases: []*ProblemTestCase{}, currentTestIndex: 0}
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
				fmt.Println(inOutFile)
				if !inOutFile.IsDir() && filepath.Ext(inOutFile.Name()) == ".in" && inFilePath == "" {
					inFilePath = fmt.Sprintf("%s/%s/%s", fullTestDirPath, dirName, inOutFile.Name())
					fmt.Println(inFilePath)
				}
				if !inOutFile.IsDir() && filepath.Ext(inOutFile.Name()) == ".ans" && outFilePath == "" {
					outFilePath = fmt.Sprintf("%s/%s/%s", fullTestDirPath, dirName, inOutFile.Name())
					fmt.Println(inFilePath)
				}
			}

			if inFilePath == "" || outFilePath == "" {
				continue
			}
			fmt.Printf("File Input Path %s\n", inFilePath)
			fmt.Printf("File Output Path %s\n", outFilePath)
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
	return nil
}
func (p *NormalProblem) NextTestCase() TestCase { //
	if len(p.TestCases) == p.currentTestIndex {
		return nil
	}
	currentTest := p.TestCases[p.currentTestIndex]
	p.currentTestIndex++
	return currentTest
}
func (p *NormalProblem) GetAllTestCases() []*ProblemTestCase {
	return p.TestCases
}

type TestCase interface {
	RunTestCase(string) (int, error)
}

type ProblemTestCase struct {
	TestInputPath  string
	TestOutputPath string
}

func NewProblemTestCase(testInputPath, testOutputPath string) *ProblemTestCase {
	return &ProblemTestCase{TestInputPath: testInputPath, TestOutputPath: testOutputPath}
}

func (t *ProblemTestCase) RunTestCase(binPath string) (int, error) {
	fmt.Println(t.TestInputPath)
	file, err := os.Open(t.TestInputPath)
	if err != nil {
		return RESULT_JUDGE_ERROR, err
	}

	cmd := exec.Command(binPath)
	cmd.Stdin = file
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	err = cmd.Run()
	if err != nil {
		return RESULT_JUDGE_ERROR, err
	}

	//	fmt.Println(t.TestOutputPath)
	return RESULT_ACCEPTED, err
}
