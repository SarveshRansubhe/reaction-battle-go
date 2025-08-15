package apis

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5/pgtype"
)

type UserJson struct {
	ID        int32              `json:"id"`
	Username  string             `json:"username"`
	Email     string             `json:"email"`
	FirstName pgtype.Text        `json:"first_name"`
	LastName  pgtype.Text        `json:"last_name"`
	CreatedAt pgtype.Timestamptz `json:"created_at"`
	UpdatedAt pgtype.Timestamptz `json:"updated_at"`
	IsActive  pgtype.Bool        `json:"is_active"`
}

type UsersResponse struct {
	UsersJson []UserJson `json:"users"`
	Count     int        `json:"count"`
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	users, err := Queries.GetAllUsers(ctx)
	if err != nil {
		log.Fatal(err)
	}

	var userJsonList []UserJson

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

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding JSON: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func CreateUsers(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	_, err := Queries.GetAllUsers(ctx)
	if err == nil {
		log.Fatal(err)
	}
}
