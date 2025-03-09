package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"golang.org/x/sys/unix"
)

var memoryLimitMb uint64
var timeLimitSec uint64
var binPath string

func main() {
	flag.StringVar(&binPath, "path", "", "Binary Path")
	flag.Uint64Var(&memoryLimitMb, "memory", 1024, "Memory Limit")
	flag.Uint64Var(&timeLimitSec, "time", 2, "Memory Limit")
	flag.Parse()
	if binPath == "" {
		log.Fatalln("No Path To binary")
	}

	if err := unix.Setrlimit(unix.RLIMIT_CPU, &unix.Rlimit{
		Cur: 2,
		Max: 2,
	}); err != nil {
		log.Fatalf("failed to set RLIMIT_CPU: %v", err)
	}

	if err := unix.Setrlimit(unix.RLIMIT_AS, &unix.Rlimit{
		Cur: 1024 * 1024 * memoryLimitMb,
		Max: 1024 * 1024 * memoryLimitMb,
	}); err != nil {
		log.Fatalf("failed to set RLIMIT_AS: %v", err)
	}
	cmd := exec.Command(binPath)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		log.Fatalf("RUNTIME ERROR")
	}

	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case err := <-done:
		log.Fatal(err)
	case <-time.After(3 * time.Second):
		fmt.Fprint(os.Stderr, "TLE")
		cmd.Process.Kill()
		fmt.Printf("TLE")
	}

}
