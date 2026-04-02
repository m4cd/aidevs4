package webserver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func middlewareCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		next.ServeHTTP(w, r)
	})
}

func CreateWebserver(handlers map[string]http.HandlerFunc, serverPort string) *http.Server {
	router := chi.NewRouter()
	for endpoint, handler := range handlers {
		router.Handle(endpoint, handler)
	}

	return &http.Server{
		Handler: middlewareCors(router),
		Addr:    ":" + serverPort,
	}
}

func RespondWithJSON(w http.ResponseWriter, code int, jsonStruct interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	respData, _ := json.Marshal(jsonStruct)
	w.Write(respData)
}

func PrintRequest(r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	fmt.Printf("Method: %s\n", r.Method)
	fmt.Printf("URL: %s\n", r.URL.String())
	fmt.Printf("Headers: %v\n", r.Header)
	fmt.Printf("Body: %s\n", string(body))
	r.Body = io.NopCloser(bytes.NewBuffer(body))
}
