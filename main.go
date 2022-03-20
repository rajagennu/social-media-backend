package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/rajagennu/social-media-backend/internal/database"
)

type errorBody struct {
	Error string `json:"error"`
}

type apiConfig struct {
	dbClient database.Client
}

func (apiCfg apiConfig) endpointUsersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		apiCfg.handlerGetUser(w, r)

	case http.MethodPost:
		apiCfg.handlerCreateUser(w, r)

	case http.MethodPut:
		apiCfg.handlerUpdateUser(w, r)

	case http.MethodDelete:
		apiCfg.handlerDeleteUser(w, r)

	default:
		respondWithError(w, 404, errors.New("method not supported"))
	}
}

func (apiCfg apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Name     string `json:"name"`
		Age      int    `json:"age"`
	}

	params := parameters{}

	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err)
		return
	}
	log.Println("docded user object ", params)

	user, err := apiCfg.dbClient.CreateUser(params.Email, params.Password, params.Name, params.Age)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err)
		return
	}
	respondWithJSON(w, http.StatusCreated, user)

}

func (apiCfg apiConfig) handlerDeleteUser(w http.ResponseWriter, r *http.Request) {
	email := strings.TrimPrefix(r.URL.Path, "/users/")
	if email == "" {
		respondWithError(w, http.StatusBadRequest, errors.New("invalid email"))
	}

	log.Println(email)
	err := apiCfg.dbClient.DeleteUser(email)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err)
		return
	}
	respondWithJSON(w, 204, "User "+email+" delete successfully.")
}

func (apiCfg apiConfig) handlerGetUser(w http.ResponseWriter, r *http.Request) {
	pathParams := r.URL.Path
	email := strings.TrimPrefix(pathParams, "/users/")
	if email == "" {
		respondWithError(w, http.StatusBadRequest, errors.New("empty email ID"))
		return
	}

	user, err := apiCfg.dbClient.GetUser(email)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err)
		return
	}

	marshalJSON, err := json.Marshal(user)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err)
	}

	respondWithJSON(w, 200, string(marshalJSON))

}

func (apiCfg apiConfig) handlerUpdateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Name     string `json:"name"`
		Age      int    `json:"age"`
	}

	params := parameters{}

	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err)
		return
	}

	email := strings.TrimPrefix(r.URL.Path, "/users/")
	if email == "" {
		respondWithError(w, http.StatusBadRequest, err)
		return
	}
	user, err := apiCfg.dbClient.UpdateUser(email, params.Password, params.Name, params.Age)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err)
		return
	}
	jsonData, _ := json.Marshal(user)
	respondWithJSON(w, http.StatusOK, string(jsonData))

}

func (apiCfg apiConfig) endPointPostsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		apiCfg.handlerGetPost(w, r)

	case http.MethodPost:
		apiCfg.handlerCreatePost(w, r)

	case http.MethodDelete:
		apiCfg.handlerDeletePost(w, r)

	}
}

func (apiCfg apiConfig) handlerGetPost(w http.ResponseWriter, r *http.Request) {

	userEmail := strings.TrimPrefix(r.URL.Path, "/posts/")

	posts, err := apiCfg.dbClient.GetPosts(userEmail)
	if err != nil {
		respondWithError(w, 404, err)
		return
	}

	data, err := json.Marshal(posts)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err)
		return
	}

	respondWithJSON(w, 200, string(data))

}

func (apiCfg apiConfig) handlerCreatePost(w http.ResponseWriter, r *http.Request) {

	type postParams struct {
		UserEmail string `json:"userEmail"`
		Text      string `json:"text"`
	}

	params := postParams{}

	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err)
		return
	}
	post, err := apiCfg.dbClient.CreatePost(params.UserEmail, params.Text)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err)
		return
	}

	respondWithJSON(w, http.StatusCreated, post)

}

func (apiCfg apiConfig) handlerDeletePost(w http.ResponseWriter, r *http.Request) {

	pathParams := r.URL.Path
	uuid := strings.TrimPrefix(pathParams, "/posts/")
	log.Printf("Received UUID %s for deletion \n", uuid)
	err := apiCfg.dbClient.DeletePost(uuid)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err)
		return
	}
	respondWithJSON(w, 204, "Post deleted successfully.")

}

func main() {

	dbClient := database.NewClient("db.json")
	err := dbClient.EnsureDB()
	if err != nil {
		log.Fatal(err)
	}

	apiCfg := apiConfig{
		dbClient: dbClient,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", testHandler)
	mux.HandleFunc("/err", testHandlerErr)

	mux.HandleFunc("/users", apiCfg.endpointUsersHandler)
	mux.HandleFunc("/users/", apiCfg.endpointUsersHandler)

	mux.HandleFunc("/posts", apiCfg.endPointPostsHandler)
	mux.HandleFunc("/posts/", apiCfg.endPointPostsHandler)

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
