package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
)

func readFileBytesByPath(pathName string) ([]byte, error) {
	file, err := os.Open(pathName)

	if err != nil {
		return nil, err
	}

	return io.ReadAll(file)
}

func removeTransimssionBytes(b []byte) []byte {
	return bytes.Map(func(r rune) rune {
		if r == 0 || r == 1 || r == 4 || r == 3 {
			return -1
		}
		return r
	}, b)
}
func CompareReaders(reader1, reader2 io.Reader) bool {
	scanner1 := bufio.NewScanner(reader1)
	scanner2 := bufio.NewScanner(reader2)
	for scanner1.Scan() && scanner2.Scan() {
		if !bytes.Equal(removeTransimssionBytes(scanner1.Bytes()), removeTransimssionBytes(scanner2.Bytes())) {
			return false
		}
	}
	return !(scanner1.Scan() && scanner2.Scan())
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	e := json.NewEncoder(w)

	err := e.Encode(v)
	if err != nil {
		return err
	}
	return nil

}
