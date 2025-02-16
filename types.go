package main

import "net/http"

type ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request) error
