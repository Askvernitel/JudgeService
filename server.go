package main

import (
	"fmt"
	"io"
	"net/http"
	_ "os"
	_ "os/exec"

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

	r.HandleFunc("/judge/{contest}/{problem}", WithErrorHandleFunc(s.getFile)).Methods("POST", "OPTIONS")
	r.Use(mux.CORSMethodMiddleware(r))
	if err := http.ListenAndServe(":4040", r); err != nil {
		return err
	}
	return nil
}

func (s *Server) getFile(w http.ResponseWriter, r *http.Request) error {
	contestName, problemName, err := getContestAndProblem(r)
	if err != nil {
		return err
	}
	err = r.ParseMultipartForm(10 << 20)
	if err != nil {
		return err
	}
	file, _, err := r.FormFile("code")

	if err != nil {
		return err
	}
	fileBytes, _ := io.ReadAll(file)

	compiler := NewCppCompiler("", &fileBytes)
	problem := NewNormalProblem(contestName, problemName)
	err = problem.initProblemTestCases()
	fmt.Println(err)
	tester := NewCppTester(problem)
	judge := NewCppJudge(compiler, tester)
	return judge.Run()
}

func getContestName(r *http.Request) (string, error) {
	contestName := mux.Vars(r)["contest"]
	if contestName == "" {
		return "", fmt.Errorf("No Contest In Path")
	}
	return contestName, nil
}

func getProblemName(r *http.Request) (string, error) {
	problemName := mux.Vars(r)["problem"]
	if problemName == "" {
		return "", fmt.Errorf("No Problem In Path")
	}
	return problemName, nil
}
func getContestAndProblem(r *http.Request) (string, string, error) {
	contestName, err := getContestName(r)
	if err != nil {
		return "", "", err
	}
	problemName, err := getProblemName(r)
	return contestName, problemName, err
}
