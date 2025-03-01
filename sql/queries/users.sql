-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_pass)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: GetUserByID :one
SELECT * FROM users WHERE 
-- name: ResetUsers :exec
DELETE FROM users;
UPDATE use
