-- name: GetAllUsers :many
SELECT * FROM public.users;

-- name: CreateUser :one
INSERT INTO public.users(
	username, email, password_hash, first_name, last_name, created_at, updated_at, is_active)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
    RETURNING *;