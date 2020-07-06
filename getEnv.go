package main

import (
	"net/http"
	"io"
	"os"
	"encoding/json"
)

func getEnv(w http.ResponseWriter, r *http.Request) {
	something, _ := json.MarshalIndent(os.Environ(), "", "   ")
	io.WriteString(w, string(something))
}