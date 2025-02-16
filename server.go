package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"

	"github.com/gorilla/mux"
)

type Server struct {
	Addr string
}

func NewServer() *Server {
	return &Server{}
}

func WithErrorHandleFunc(f ErrorHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func (s *Server) Run() error {
	r := mux.NewRouter()

	r.HandleFunc("/file/upload", WithErrorHandleFunc(s.getFile)).Methods("POST", "OPTIONS")
	r.Use(mux.CORSMethodMiddleware(r))
	if err := http.ListenAndServe(":4040", r); err != nil {
		return err
	}
	return nil
}

func (s *Server) getFile(w http.ResponseWriter, r *http.Request) error {
	err := r.ParseMultipartForm(10 << 20)
	fmt.Println("Hi")
	if err != nil {
		return err
	}
	file, _, err := r.FormFile("code")
	if err != nil {
		return err
	}
	fileBytes, _ := io.ReadAll(file)

	tmpFile, err := os.CreateTemp("./uploaded-files-tmp/", "file*.cpp")
	if err != nil {
		return err
	}

	_, err = tmpFile.Write(fileBytes)
	if err != nil {
		return err
	}
	input := "50\n50\n"
	fmt.Println(tmpFile.Name())
	cmd := exec.Command("g++", tmpFile.Name(), "-o", fmt.Sprintf("./uploaded-files-tmp/output"))

	err = cmd.Run()
	if err != nil {
		return err
	}

	cmd = exec.Command("./uploaded-files-tmp/output")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	_, err = io.WriteString(stdin, input)
	if err != nil {
		return err
	}

	//	cmd.Stdout = os.Stdout
	output, err := cmd.Output()
	if err != nil {
		return err
	}
	fmt.Println(string(output))
	return nil
}
