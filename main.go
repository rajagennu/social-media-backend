package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/rajagennu/social-media-backend/internal/database"
)

type errorBody struct {
	Error string `json:"error"`
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", testHandler)
	mux.HandleFunc("/err", testHandlerErr)

	const addr = "localhost:8080"
	srv := http.Server{
		Handler:      mux,
		Addr:         addr,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  30 * time.Second,
	}
	srv.ListenAndServe()
}

// func testHandler(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(200)
// 	w.Write([]byte("{}"))
// }

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func respondWithError(w http.ResponseWriter, code int, err error) {

	e := errorBody{
		Error: err.Error(),
	}

	respondWithJSON(w, 404, e)
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	// you can use any compatible type, but let's use our database package's User type for practice
	respondWithJSON(w, 200, database.User{
		Email: "test@example.com",
	})
}

func testHandlerErr(w http.ResponseWriter, r *http.Request) {
	// you can use any compatible type, but let's use our database package's User type for practice
	respondWithError(w, 404, errors.New("404 not found"))
}
