package main

import (
	//	"bytes"
	"fmt"
	"io"
	"os/exec"

	"github.com/google/uuid"
)

type Compiler interface {
	Compile() error
	OutputBinPathName() string
}

const (
	GPP_COMPILER_COMMAND  = "g++"
	COMPILED_BINARIES_DIR = "./uploaded-files-tmp/"
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
		OutputFileName: COMPILED_BINARIES_DIR + uuid.New().String(),
	}
}

func (c *GppCompiler) Compile() error {
	if c.FilePathName == "" {
		//TODO: make this readable command
		fmt.Println(c.OutputFileName)
		cmd := exec.Command(GPP_COMPILER_COMMAND, "-o", c.OutputFileName, "-xc++", "-")
		stdin, err := cmd.StdinPipe()
		if err != nil {
			return fmt.Errorf("failed to get stdin pipe: %w", err)
		}
		if err := cmd.Start(); err != nil {
			return fmt.Errorf("failed to start command: %w", err)
		}
		_, err = io.WriteString(stdin, string(*c.FileData))
		if err != nil {
			return fmt.Errorf("failed to write to stdin: %w", err)
		}
		stdin.Close()
		return cmd.Wait()
	}
	cmd := exec.Command(GPP_COMPILER_COMMAND, c.FilePathName, "-o", c.OutputFileName, "-xc++ -")

	return cmd.Run()
}

func (c *GppCompiler) OutputBinPathName() string {
	return c.OutputFileName
}

func (c *GppCompiler) DeleteOutputFile() {
}
