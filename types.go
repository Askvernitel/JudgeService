package main

import (
	"net/http"
)

type ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request) error

type TestLimits struct {
	MemoryLimitMb int64 `yaml:"memory_limit_mb" json:"memoryLimitMb"`
	TimeLimitSec  int64 `yaml:"time_limit_sec" json:"timeLimitMb"`
}

type ProblemYaml struct {
	TestsPath   string `yaml:"tests_path"`
	*TestLimits `yaml:"limits"`
}
type JudgeResponse struct {
	Results         []int `json:"results"`
	TimeForEachTest []int `json:"timeForEachTest"`
	Score           int   `json:"score"`
}
type ApiError struct {
	Error string `json:"error"`
}

type CmdResult struct {
	Result        int
	TimeTakenSec  int
	MemoryTakenMb int
}
