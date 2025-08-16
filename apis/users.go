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

func UserApi(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		GetUsers(w, r)

	case http.MethodPost:
		CreateUser(w, r)

	default:
		log.Printf("Method not implemented")
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	users, err := Queries.GetAllUsers(ctx)
	if err != nil {
		log.Println(err)
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
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	var arg UserJson

	err := json.NewDecoder(r.Body).Decode(&arg)
	if err != nil {
		log.Println(err)
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
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Validations
	// if(u.Username == ni)

	returnUser := UserJson{
		Username:  u.Username,
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		IsActive:  u.IsActive,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(returnUser); err != nil {
		log.Printf("Error encoding JSON: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func GetPgTime(time time.Time) pgtype.Timestamptz {
	return pgtype.Timestamptz{
		Time:  time,
		Valid: true,
	}
}
