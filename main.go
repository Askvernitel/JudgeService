package main

//TODO: Just Leaving These maybe add buffering to files so all
// data will not just go to memory

//TODO: Add every path variable to path
import (
	"log"
	"os"
)

var problemsPath string

func main() {
	problemsPath = os.Getenv("PROBLEMS_PATH")
	//binOutPath := os.Getenv("BIN_OUTPUT_PATH")

	httpServerPort := os.Getenv("SERVER_PORT")

	server := NewServer(httpServerPort)

	err := server.Run()
	if err != nil {
		log.Fatal("Server Did Not Run")

	}
}
