package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type User struct {
	ID           string
	Email        string
	Username     string
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type DB struct{ pool *pgxpool.Pool }

func (db *DB) Pool() *pgxpool.Pool {
	return db.pool
}

func NewDB(ctx context.Context, connString string) (*DB, error) {
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("new pool: %w", err)
	}
	return &DB{pool: pool}, nil
}

func (db *DB) CreateUser(ctx context.Context, u User) error {

	_, err := db.pool.Exec(
		ctx,
		`
		INSERT INTO users (id, email, username, password_hash)
		VALUES ($1, $2, $3, $4)
		`,
		u.ID,
		u.Email,
		u.Username,
		u.PasswordHash,
	)

	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}

	return nil
}

func (db *DB) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	u := &User{}
	err := db.pool.QueryRow(
		ctx,
		`SELECT id, email, username, password_hash, created_at, updated_at FROM users WHERE email = $1`,
		email,
	).Scan(&u.ID, &u.Email, &u.Username, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("get user by email: %w", err)
	}
	return u, nil
}

func (db *DB) GetUserByID(ctx context.Context, id string) (*User, error) {
	u := &User{}
	err := db.pool.QueryRow(
		ctx,
		`SELECT id, email, username, password_hash, created_at, updated_at FROM users WHERE id = $1`,
		id,
	).Scan(&u.ID, &u.Email, &u.Username, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("get user by id: %w", err)
	}
	return u, nil
}

func (db *DB) GetUserByUsername(ctx context.Context, username string) (*User, error) {
	u := &User{}
	err := db.pool.QueryRow(
		ctx,
		`SELECT id, email, username, password_hash, created_at, updated_at FROM users WHERE username = $1`,
		username,
	).Scan(&u.ID, &u.Email, &u.Username, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("get user by username: %w", err)
	}
	return u, nil
}

func (db *DB) UpdateUser(ctx context.Context, u User) error {
	_, err := db.pool.Exec(
		ctx,
		`
		UPDATE users
		SET email = $1, username = $2, password_hash = $3
		WHERE id = $4
		`,
		u.Email,
		u.Username,
		u.PasswordHash,
		u.ID,
	)
	if err != nil {
		return fmt.Errorf("update user: %w", err)
	}
	return nil
}

func (db *DB) DeleteUser(ctx context.Context, id string) error {
	_, err := db.pool.Exec(ctx, `DELETE FROM users WHERE id = $1`, id)

	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}
	return nil
}
