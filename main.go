package main

//TODO: Just Leaving These maybe add buffering to files so all
// data will not jut go to memory
import "os"

var problemsPath string

func main() {
	os.Setenv("PROBLEMS_PATH", "./problems")
	os.Setenv("BIN_OUTPUT_PATH", "~/Desktop/backend-project/JudgeService/uploaded-files-tmp")
	problemsPath = os.Getenv("PROBLEMS_PATH")

	server := NewServer()

	server.Run()
}
