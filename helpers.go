package main

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"os"
)

func readFileBytesByPath(pathName string) ([]byte, error) {
	file, err := os.Open(pathName)

	if err != nil {
		return nil, err
	}

	return io.ReadAll(file)
}
func CompareReaders(reader1, reader2 io.Reader) bool {
	scanner1 := bufio.NewScanner(reader1)
	scanner2 := bufio.NewScanner(reader2)
	for scanner1.Scan() && scanner2.Scan() {
		log.Println(scanner1.Text())
		if !bytes.Equal(scanner1.Bytes(), scanner2.Bytes()) {
			return false
		}
	}

	return !(scanner1.Scan() && scanner2.Scan())
}
