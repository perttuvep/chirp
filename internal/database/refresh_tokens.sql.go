// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: refresh_tokens.sql

package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

const createRToken = `-- name: CreateRToken :one
INSERT INTO refresh_tokens (token, created_at, updated_at, user_id, expires_at,revoked_at)
VALUES (
    $1,
    NOW(),
    NOW(),
    $2,
    $3,
    $4
)
RETURNING token, created_at, updated_at, user_id, expires_at, revoked_at
`

type CreateRTokenParams struct {
	Token     string
	UserID    uuid.UUID
	ExpiresAt time.Time
	RevokedAt sql.NullTime
}

func (q *Queries) CreateRToken(ctx context.Context, arg CreateRTokenParams) (RefreshToken, error) {
	row := q.db.QueryRowContext(ctx, createRToken,
		arg.Token,
		arg.UserID,
		arg.ExpiresAt,
		arg.RevokedAt,
	)
	var i RefreshToken
	err := row.Scan(
		&i.Token,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.UserID,
		&i.ExpiresAt,
		&i.RevokedAt,
	)
	return i, err
}

const getRTokenByToken = `-- name: GetRTokenByToken :one
SELECT token, created_at, updated_at, user_id, expires_at, revoked_at FROM refresh_tokens WHERE token = $1
`

func (q *Queries) GetRTokenByToken(ctx context.Context, token string) (RefreshToken, error) {
	row := q.db.QueryRowContext(ctx, getRTokenByToken, token)
	var i RefreshToken
	err := row.Scan(
		&i.Token,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.UserID,
		&i.ExpiresAt,
		&i.RevokedAt,
	)
	return i, err
}

const getRTokenByUID = `-- name: GetRTokenByUID :one
SELECT token, created_at, updated_at, user_id, expires_at, revoked_at FROM refresh_tokens WHERE user_id = $1
`

func (q *Queries) GetRTokenByUID(ctx context.Context, userID uuid.UUID) (RefreshToken, error) {
	row := q.db.QueryRowContext(ctx, getRTokenByUID, userID)
	var i RefreshToken
	err := row.Scan(
		&i.Token,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.UserID,
		&i.ExpiresAt,
		&i.RevokedAt,
	)
	return i, err
}

const getUserFromRTkon = `-- name: GetUserFromRTkon :one
SELECT id, users.created_at, users.updated_at, email, hashed_pass, is_chirpy_red, token, refresh_tokens.created_at, refresh_tokens.updated_at, user_id, expires_at, revoked_at FROM users INNER JOIN refresh_tokens on refresh_tokens(user_id) = users(id) WHERE refresh_tokens(token) = $1
`

type GetUserFromRTkonRow struct {
	ID          uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Email       string
	HashedPass  string
	IsChirpyRed bool
	Token       string
	CreatedAt_2 time.Time
	UpdatedAt_2 time.Time
	UserID      uuid.UUID
	ExpiresAt   time.Time
	RevokedAt   sql.NullTime
}

func (q *Queries) GetUserFromRTkon(ctx context.Context, token string) (GetUserFromRTkonRow, error) {
	row := q.db.QueryRowContext(ctx, getUserFromRTkon, token)
	var i GetUserFromRTkonRow
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Email,
		&i.HashedPass,
		&i.IsChirpyRed,
		&i.Token,
		&i.CreatedAt_2,
		&i.UpdatedAt_2,
		&i.UserID,
		&i.ExpiresAt,
		&i.RevokedAt,
	)
	return i, err
}

const revokeToken = `-- name: RevokeToken :exec
UPDATE refresh_tokens
    SET revoked_at = NOW(),updated_at = NOW()
    WHERE token = $1
`

func (q *Queries) RevokeToken(ctx context.Context, token string) error {
	_, err := q.db.ExecContext(ctx, revokeToken, token)
	return err
}
