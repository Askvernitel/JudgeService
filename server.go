package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	_ "os"
	_ "os/exec"

	"github.com/gorilla/handlers"
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

	r.Use(mux.CORSMethodMiddleware(r))
	r.HandleFunc("/judge/{contest}/{problem}", WithErrorHandleFunc(s.JudgeProblem)).Methods("POST", "OPTIONS")

	if err := http.ListenAndServe(":4040", handlers.CORS()(r)); err != nil {
		return err
	}
	return nil
}

func (s *Server) JudgeProblem(w http.ResponseWriter, r *http.Request) error {
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
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return err
	}
	compiler := NewCppCompiler("", &fileBytes)
	problem := NewNormalProblem(contestName, problemName)
	err = problem.initProblemTestCases()
	if err != nil {
		return err
	}
	tester := NewCppTester(problem)
	judge := NewCppJudge(compiler, tester)

	err = judge.Run()
	if err != nil {
		return err
	}
	fmt.Println(judge.Results)
	err = json.NewEncoder(w).Encode(judge.Results)
	if err != nil {
		return err
	}
	return nil
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
