package apis

import (
	"app/sql/datastore"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type UserJson struct {
	ID           int32              `json:"id"`
	Username     string             `json:"username"`
	Email        string             `json:"email"`
	PasswordHash string             `json:"password_hash,omitempty"`
	FirstName    pgtype.Text        `json:"first_name"`
	LastName     pgtype.Text        `json:"last_name"`
	CreatedAt    pgtype.Timestamptz `json:"created_at"`
	UpdatedAt    pgtype.Timestamptz `json:"updated_at"`
	IsActive     pgtype.Bool        `json:"is_active"`
}

type UsersResponse struct {
	UsersJson []UserJson `json:"users"`
	Count     int        `json:"count"`
}

type ErrorResponse struct {
	Message    string `json:"message"`
	Stacktrace string `json:"stacktrace"`
}

func UserApi(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		GetUsers(w, r)
	case http.MethodPost:
		CreateUser(w, r)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	users, err := Queries.GetAllUsers(ctx)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	var userJsonList []UserJson = []UserJson{}
	for _, u := range users {
		r := UserJson{
			ID:        u.ID,
			Username:  u.Username,
			Email:     u.Email,
			FirstName: u.FirstName,
			LastName:  u.LastName,
			CreatedAt: u.CreatedAt,
			UpdatedAt: u.UpdatedAt,
			IsActive:  u.IsActive,
		}
		userJsonList = append(userJsonList, r)
	}
	response := UsersResponse{
		UsersJson: userJsonList,
		Count:     len(users),
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding JSON: %v", err)
	}
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	var arg UserJson
	err := json.NewDecoder(r.Body).Decode(&arg)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	if usernameDup, err := Queries.CheckDuplicateUsername(ctx, arg.Username); err != nil {
		log.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	} else if usernameDup > 0 {
		http.Error(w, "Username already exists", http.StatusConflict)
		return
	}
	if emailDup, err := Queries.CheckDuplicateEmail(ctx, arg.Email); err != nil {
		log.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	} else if emailDup > 0 {
		http.Error(w, "Email already exists", http.StatusConflict)
		return
	}
	u, err := Queries.CreateUser(ctx, datastore.CreateUserParams{
		Username:     arg.Username,
		Email:        arg.Email,
		PasswordHash: arg.PasswordHash,
		FirstName:    arg.FirstName,
		LastName:     arg.LastName,
		CreatedAt:    GetPgTime(time.Now()),
		UpdatedAt:    GetPgTime(time.Now()),
		IsActive:     pgtype.Bool{Bool: true, Valid: true},
	})
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	returnUser := UserJson{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		IsActive:  u.IsActive,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(returnUser); err != nil {
		log.Printf("Error encoding JSON: %v", err)
	}
}

func GetPgTime(time time.Time) pgtype.Timestamptz {
	return pgtype.Timestamptz{
		Time:  time,
		Valid: true,
	}
}
