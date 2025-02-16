package main

import (
	//	"bytes"
	"fmt"
	"os"
	"os/exec"
)

type Compiler interface {
	Compile() error
}

const (
	GPP_COMPILER_COMMAND = "g++"
)

type GppCompiler struct {
	FilePathName   string
	FileData       *[]byte
	OutputFileName string
}

func NewCppCompiler(filePathName string, fileData *[]byte) *GppCompiler {
	return &GppCompiler{
		FilePathName:   filePathName,
		FileData:       fileData,
		OutputFileName: "Output",
	}
}

func (c *GppCompiler) Compile() error {
	if c.FilePathName == "" {
		//TODO: make this readable command
		fmt.Println(c.OutputFileName)
		//		fmt.Println(string(*c.FileData))
		cmd := exec.Command("bash", "-c", fmt.Sprintf("echo 'echo %s | %s %s %s %s'", string(*c.FileData), GPP_COMPILER_COMMAND, "-o", c.OutputFileName, "-xc++ -"))

		//	dataBuf := bytes.NewBuffer(*c.FileData)

		//	cmd.Stdin = dataBuf
		cmd.Stderr = os.Stdout //for debug
		err := cmd.Run()
		if err != nil {
			return err
		}

		return nil
	}
	cmd := exec.Command(GPP_COMPILER_COMMAND, c.FilePathName, "-o", c.OutputFileName, "-xc++ -")
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func (c *GppCompiler) DeleteOutputFile() {
}
