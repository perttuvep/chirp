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
SELECT * FROM users WHERE id = $1;
-- name: ResetUsers :exec
DELETE FROM users;

-- name: EditUserEmail :one
UPDATE users
    SET email = $2, hashed_pass = $3, updated_at = NOW()
    WHERE id = $1
RETURNING *;

-- name: ChirpyRedEnableByID :exec
UPDATE users
    SET is_chirpy_red = true, updated_at = NOW()
WHERE id = $1;
-- name: ChirpyRedDisableByID :exec
UPDATE users
    SET is_chirpy_red = false, updated_at = NOW()
WHERE id = $1;

