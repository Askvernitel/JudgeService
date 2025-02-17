package main

import (
	//	"bytes"
	"fmt"
	"os"
	"os/exec"
)

type Compiler interface {
	Compile() error
	OutputBinPathName() string
}

const (
	GPP_COMPILER_COMMAND = "g++"
)

type GppCompiler struct {
	FilePathName   string
	FileData       *[]byte
	OutputFileName string
}

// TODO: Give to constructor output path
func NewCppCompiler(filePathName string, fileData *[]byte) *GppCompiler {
	return &GppCompiler{
		FilePathName: filePathName, FileData: fileData,
		OutputFileName: "./uploaded-files-tmp/output",
	}
}

func (c *GppCompiler) Compile() error {
	if c.FilePathName == "" {
		//TODO: make this readable command
		fmt.Println(c.OutputFileName)
		cmd := exec.Command("sh", "-c", fmt.Sprintf("echo '%s' | %s %s %s %s", string(*c.FileData), GPP_COMPILER_COMMAND, "-o", c.OutputFileName, "-xc++ -"))

		cmd.Stderr = os.Stdout

		return cmd.Run()
	}
	cmd := exec.Command(GPP_COMPILER_COMMAND, c.FilePathName, "-o", c.OutputFileName, "-xc++ -")

	return cmd.Run()
}

func (c *GppCompiler) OutputBinPathName() string {
	return c.OutputFileName
}

func (c *GppCompiler) DeleteOutputFile() {
}
