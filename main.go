package main

import "os"

func main() {
	os.Setenv("PROBLEMS_PATH", "~/Desktop/backend-project/JudgeService/problems/")
	os.Setenv("BIN_OUTPUT_PATH", "~/Desktop/backend-project/JudgeService/uploaded-files-tmp/")
	server := NewServer()

	server.Run()
}
