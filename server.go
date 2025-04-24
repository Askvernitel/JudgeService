package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	_ "os/exec"

	pb "JudgeService.com/proto"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	UPLOADED_CODE_FORM_KEY = "code"
	MAX_UPLOAD_SIZE        = 1024 * 1024 * 5
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
			WriteJSON(w, http.StatusInternalServerError, ApiError{Error: err.Error()})
			//			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func WithAuthHandleFunc(f http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		token, err := ExtractToken(r)
		if err != nil {
			WriteJSON(w, http.StatusForbidden, ApiError{Error: err.Error()})
			return
		}
		conn, err := grpc.NewClient("localhost:50000", grpc.WithTransportCredentials(insecure.NewCredentials()))

		if err != nil {
			WriteJSON(w, http.StatusInternalServerError, ApiError{Error: err.Error()})
			return
		}
		defer conn.Close()
		service := pb.NewAuthServiceClient(conn)
		resp, err := service.Auth(context.Background(), &pb.AuthRequest{Token: token})

		if err != nil {
			WriteJSON(w, http.StatusInternalServerError, ApiError{Error: err.Error()})
			return
		}
		fmt.Println(resp.Ok)

		f(w, r)
	}
}

func (s *Server) Run() error {
	r := mux.NewRouter()

	r.Use(mux.CORSMethodMiddleware(r))
	r.HandleFunc("/judge/{contest}/{problem}", WithAuthHandleFunc(WithErrorHandleFunc(s.JudgeProblem))).Methods("POST", "OPTIONS")

	if err := http.ListenAndServe(":4040", handlers.CORS()(r)); err != nil {
		return err
	}

	return nil

}

func getFileData(r *http.Request, key string) ([]byte, error) {
	err := r.ParseMultipartForm(MAX_UPLOAD_SIZE)
	if err != nil {
		return nil, err
	}
	file, _, err := r.FormFile(key)
	if err != nil {
		return nil, err
	}
	return io.ReadAll(file)
}
func RunJudge(w http.ResponseWriter, fileBytes []byte) error {

	return nil
}

func (s *Server) JudgeProblem(w http.ResponseWriter, r *http.Request) error {
	r.Body = http.MaxBytesReader(w, r.Body, MAX_UPLOAD_SIZE)
	contestName, problemName, err := getContestAndProblem(r)
	if err != nil {
		return err
	}
	fileBytes, err := getFileData(r, UPLOADED_CODE_FORM_KEY)
	if err != nil {
		return err
	}
	compiler := NewCppCompiler("", &fileBytes)
	defer func() {
		err := compiler.DeleteOutputFile()
		if err != nil {
			panic(err)
		}
	}()
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
	return WriteJSON(w, http.StatusOK, JudgeResponse{TestResults: judge.Results})
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
