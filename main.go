package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"main/middleware"
	"main/users"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/lib/pq"
)

type ConnectionDAO struct {
	queries *users.Queries
}

type UserInput struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	UserType    string `json:"user_type"`
	Address     string `json:"address"`
}

type UserResponse struct {
	ID          int32  `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	UserType    string `json:"user_type"`
	Address     string `json:"address"`
}

func main() {
	connStr := os.Getenv("DATABASE_URI")
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	queries := users.New(db)
	dao := &ConnectionDAO{queries: queries}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /users", dao.handleGetAllUsers)

	mux.HandleFunc("GET /users/{id}", dao.handleGetUser)

	mux.HandleFunc("PUT /users/{id}", dao.handleUpdateUser)

	mux.HandleFunc("DELETE /users/{id}", dao.handleDeleteUser)

	mux.HandleFunc("POST /users", dao.handleCreateUser)

	server := &http.Server{
		Addr:    ":8080",
		Handler: middleware.Logging(mux),
	}

	log.Println("Server started, listening on port 8080")
	server.ListenAndServe()

	defer db.Close()
}

func (dao *ConnectionDAO) handleGetAllUsers(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	users, err := dao.queries.ListUsers(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var userResponseArray []UserResponse
	for _, user := range users {
		userResponseArray = append(userResponseArray, UserResponse{
			ID:          user.ID,
			Name:        user.Name,
			Email:       user.Email,
			PhoneNumber: user.PhoneNumber,
			UserType:    string(user.UserType),
			Address:     user.Address,
		})
	}

	if userResponseArray == nil {
		userResponseArray = []UserResponse{}
	}

	jsonResponse, err := json.Marshal(userResponseArray)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func (dao *ConnectionDAO) handleGetUser(w http.ResponseWriter, r *http.Request) {
	identifier := r.PathValue("id")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user users.UserAccount
	var err error

	// Check if the identifier is an integer (ID)
	id, err := strconv.ParseInt(identifier, 10, 32)
	if err == nil {
		// If the identifier is an integer, use GetUserById query
		user, err = dao.queries.GetUserById(ctx, int32(id))
	} else {
		// If the identifier is not an integer, assume it's an email and use GetUserByEmail query
		user, err = dao.queries.GetUserByEmail(ctx, identifier)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	var userResponse UserResponse
	userResponse.ID = user.ID
	userResponse.Name = user.Name
	userResponse.Email = user.Email
	userResponse.PhoneNumber = user.PhoneNumber
	userResponse.UserType = string(user.UserType)
	userResponse.Address = user.Address

	jsonResponse, err := json.Marshal(userResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func (dao *ConnectionDAO) handleUpdateUser(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	idInt, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var input UserInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user := users.UpdateUserParams{
		ID:          int32(idInt),
		Name:        input.Name,
		Email:       input.Email,
		PhoneNumber: input.PhoneNumber,
		UserType:    users.UserType(input.UserType),
		Address:     input.Address,
	}

	err = dao.queries.UpdateUser(ctx, user)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (dao *ConnectionDAO) handleDeleteUser(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	idInt, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = dao.queries.DeleteUser(ctx, int32(idInt))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (dao *ConnectionDAO) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.NotFound(w, r)
		return
	}

	var input UserInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check if email is already in use
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	foundUser, err := dao.queries.GetUserByEmail(ctx, input.Email)
	if err == nil && foundUser.Email == input.Email {
		http.Error(w, "Email already in use", http.StatusBadRequest)
		return
	}

	user := users.CreateUserParams{
		Name:        input.Name,
		Email:       input.Email,
		PhoneNumber: input.PhoneNumber,
		UserType:    users.UserType(input.UserType),
		Address:     input.Address,
	}

	createdUser, err := dao.queries.CreateUser(ctx, user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var userResponse UserResponse
	userResponse.ID = createdUser.ID
	userResponse.Name = createdUser.Name
	userResponse.Email = createdUser.Email
	userResponse.PhoneNumber = createdUser.PhoneNumber
	userResponse.UserType = string(createdUser.UserType)
	userResponse.Address = createdUser.Address

	jsonResponse, err := json.Marshal(userResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonResponse)
}
