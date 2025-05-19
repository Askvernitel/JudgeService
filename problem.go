package main

import (
	"bytes"
	"fmt"
	"io"
	_ "io"
	"os"
	_ "os/exec"
	"path/filepath"

	yaml "gopkg.in/yaml.v3"
)

// PROBLEM
type Problem interface {
	NextTestCase() TestCase
	GetAllTestCases() []*ProblemTestCase
	GetTestLimits() *TestLimits
}

const (
	PROBLEM_YAML            = "problem.yaml"
	DEFAULT_TEST_PATH       = "data"
	DEFAULT_MEMORY_LIMIT_MB = 256
	DEFAULT_TIME_LIMIT_SEC  = 1
)

type NormalProblem struct {
	ProblemPathName  string
	TestsDirPathName string
	TestCases        []*ProblemTestCase
	TestLimits       *TestLimits
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
	problemInfo := &ProblemYaml{TestsPath: DEFAULT_TEST_PATH, TestLimits: &TestLimits{MemoryLimitMb: DEFAULT_MEMORY_LIMIT_MB, TimeLimitSec: DEFAULT_TIME_LIMIT_SEC}}
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
			p.TestCases = append(p.TestCases, NewProblemTestCase(inFilePath, outFilePath, p.TestLimits))
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
	p.TestLimits = problemInfo.TestLimits
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
func (p *NormalProblem) GetTestLimits() *TestLimits {
	return p.TestLimits
}

// TESTCASE
type TestCase interface {
	RunTestCase(string, ResourceLimiter) (*TestResult, error)
}

type ProblemTestCase struct {
	TestInputPath  string
	TestOutputPath string
	TestLimits     *TestLimits
}

func NewProblemTestCase(testInputPath, testOutputPath string, testLimits *TestLimits) *ProblemTestCase {
	return &ProblemTestCase{TestInputPath: testInputPath, TestOutputPath: testOutputPath, TestLimits: testLimits}
}

func (t *ProblemTestCase) isCorrectOutput(out io.Reader) (int, error) {
	//fmt.Println(t.TestOutputPath)
	correctOutputFile, err := os.Open(t.TestOutputPath)
	if err != nil {
		return RESULT_JUDGE_ERROR, err
	}
	defer correctOutputFile.Close()
	areEqualReaders := CompareReaders(correctOutputFile, out)
	if !areEqualReaders {
		return RESULT_WRONG_ANSWER, nil
	}
	return RESULT_ACCEPTED, nil

}
func (t *ProblemTestCase) RunTestCase(binPath string, cmd ResourceLimiter) (*TestResult, error) {
	//	cmd := NewCmdLimiter(binPath, t.TestLimits.MemoryLimitMb, t.TestLimits.TimeLimitSec)
	testResult := &TestResult{}
	inputFile, err := os.Open(t.TestInputPath)
	if err != nil {
		return WriteResult(testResult, RESULT_JUDGE_ERROR, 0), err
	}
	defer inputFile.Close()
	var clientOutputBuffer bytes.Buffer

	cmd.SetStdin(inputFile)
	cmd.SetStdout(&clientOutputBuffer)
	cmdResult, err := cmd.Run()
	if err != nil {
		return WriteResult(testResult, RESULT_JUDGE_ERROR, 0), err
	}
	//choose the result
	switch cmdResult.Result {
	case CMD_RESULT_RUN_SUCCESSFUL:
		result, err := t.isCorrectOutput(&clientOutputBuffer)
		if err != nil {
			return WriteResult(testResult, RESULT_JUDGE_ERROR, 0), err
		}
		return WriteResult(testResult, result, cmdResult.TimeTakenSec), nil
	case CMD_RESULT_TIME_EXCEEDED_LIMIT:
		return WriteResult(testResult, RESULT_TIME_EXCEEDED_LIMIT, 0), nil
	case CMD_RESULT_MEMORY_EXCEEDED_LIMIT:
		return WriteResult(testResult, RESULT_MEMORY_EXCEEDED_LIMIT, 0), nil
	default:
		return WriteResult(testResult, RESULT_JUDGE_ERROR, 0), nil
	}
}
