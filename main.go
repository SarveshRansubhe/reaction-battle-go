package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"

	"app/apis"
	"app/sql/datastore"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, Go HTTP Server!")
}

func connectPostgres() {
	connString := os.Getenv("POSTGRES_CONNECTION_STRING")
	fmt.Println("Connection String:", connString)

	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Connected to db", conn.Config().Database)
	apis.Queries = datastore.New(conn)
	// defer conn.Close(context.Background())
	// defer fmt.Println("Disconnected from Postgres")
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/getUsers", apis.GetUsers)
	http.HandleFunc("/*", http.NotFound)

	connectPostgres()

	fmt.Println("Server is running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
