-- name: CheckDuplicateUsername :one
SELECT Count(*) from public.users where username = $1;

-- name: CheckDuplicateEmail :one
SELECT Count(*) from public.users where email = $1;

-- name: CreateUser :one
INSERT INTO public.users(
	username, email, password_hash, first_name, last_name, created_at, updated_at, is_active)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
    RETURNING *;

-- name: GetAllUsers :many
SELECT * FROM public.users;