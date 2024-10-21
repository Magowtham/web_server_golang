package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
)

// user model
type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// inmemory database
var database map[int]User = make(map[int]User)

// mutex
var dbMutex sync.RWMutex

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s", "Hello Developer 💀 👋...")
}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	error := json.NewDecoder(r.Body).Decode(&user)

	if error != nil {
		http.Error(w, "invalid json format", http.StatusBadRequest)
		return
	}

	if user.Name == "" {
		http.Error(w, "user name required", http.StatusBadRequest)
		return
	}

	if user.Email == "" {
		http.Error(w, "email required", http.StatusBadRequest)
		return
	}

	dbMutex.Lock()
	defer dbMutex.Unlock()

	for _, value := range database {
		if value.Email == user.Email {
			http.Error(w, "email already exists", http.StatusBadRequest)
			return
		}
	}

	database[len(database)+1] = user

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "user created successfully")
}

func getUserByIdHandler(w http.ResponseWriter, r *http.Request) {
	id, error := strconv.Atoi(r.PathValue("id"))

	if error != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	dbMutex.RLock()
	defer dbMutex.RUnlock()
	user := database[id]

	if user.Name == "" {
		http.Error(w, "user not found", http.StatusBadRequest)
		return
	}

	res, error := json.Marshal(user)

	if error != nil {
		http.Error(w, "error occurred while encoding to json", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func deleteUserByIdHandler(w http.ResponseWriter, r *http.Request) {
	id, error := strconv.Atoi(r.PathValue("id"))

	if error != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	dbMutex.Lock()
	defer dbMutex.Unlock()

	user := database[id]

	if user.Name == "" {
		http.Error(w, "user not found", http.StatusBadRequest)
		return
	}

	delete(database, id)

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "user deleted successfully")

}

func main() {
	mux := http.NewServeMux()

	//routes
	mux.HandleFunc("GET /", rootHandler)
	mux.HandleFunc("POST /user", createUserHandler)
	mux.HandleFunc("GET /user/{id}", getUserByIdHandler)
	mux.HandleFunc("DELETE /user/{id}", deleteUserByIdHandler)

	fmt.Printf("server is listening on 0.0.0.0:8000\n")
	http.ListenAndServe("0.0.0.0:8000", mux)
}
