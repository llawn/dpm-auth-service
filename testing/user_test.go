package auth_test

import (
	"context"
	"errors"
	"testing"

	"dpm/services/auth"

	"github.com/jackc/pgx/v5"
	"github.com/google/go-cmp/cmp"
)

func mockUser() auth.User {
	return auth.User{
		ID:           "5b6ad7b1-dd52-412d-bc1f-e85e871311a5",
		Email:        "test@test.com",
		Username:     "testusername",
		PasswordHash: "$123abc.",
	}
}

func createUser(t *testing.T, db *auth.DB, ctx context.Context, user auth.User) {
	t.Helper()
	err := db.CreateUser(ctx, user)
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}
}

func getUserByEmail(t *testing.T, db *auth.DB, ctx context.Context, want auth.User) {
	t.Helper()
	got, err := db.GetUserByEmail(ctx, want.Email)
	if err != nil {
		t.Fatalf("failed to get user by email: %v", err)
	}
	want.CreatedAt = got.CreatedAt
	want.UpdatedAt = got.UpdatedAt
	if diff := cmp.Diff(&want, got); diff != "" {
		t.Errorf("GetUserByEmail() mismatch (-want +got):\n%s", diff)
	}
}

func getUserByID(t *testing.T, db *auth.DB, ctx context.Context, want auth.User) {
	t.Helper()
	got, err := db.GetUserByID(ctx, want.ID)
	if err != nil {
		t.Fatalf("failed to get user by id: %v", err)
	}
	want.CreatedAt = got.CreatedAt
	want.UpdatedAt = got.UpdatedAt
	if diff := cmp.Diff(&want, got); diff != "" {
		t.Errorf("GetUserByID() mismatch (-want +got):\n%s", diff)
	}
}

func getUserByUsername(t *testing.T, db *auth.DB, ctx context.Context, want auth.User) {
	t.Helper()
	got, err := db.GetUserByUsername(ctx, want.Username)
	if err != nil {
		t.Fatalf("failed to get user by username: %v", err)
	}
	want.CreatedAt = got.CreatedAt
	want.UpdatedAt = got.UpdatedAt
	if diff := cmp.Diff(&want, got); diff != "" {
		t.Errorf("GetUserByUsername() mismatch (-want +got):\n%s", diff)
	}
}

func updateUser(t *testing.T, db *auth.DB, ctx context.Context, user auth.User) auth.User {
	t.Helper()
	user.Username = "alice"
	err := db.UpdateUser(ctx, user)
	if err != nil {
		t.Fatalf("failed to update user: %v", err)
	}
	updated, err := db.GetUserByID(ctx, user.ID)
	if err != nil {
		t.Fatalf("failed to fetch updated user: %v", err)
	}
	if updated.Username != user.Username {
		t.Errorf("expected username to be %s, got %s", user.Username, updated.Username)
	}
	return *updated
}

func deleteUser(t *testing.T, db *auth.DB, ctx context.Context, id string) {
	t.Helper()
	del_err := db.DeleteUser(ctx, id)
	if del_err != nil {
		t.Fatalf("failed to delete user: %v", del_err)
	}
	_, get_err := db.GetUserByID(ctx, id)
	if !errors.Is(get_err, pgx.ErrNoRows) {
		t.Errorf("expected pgx.ErrNoRows, got %v", get_err)
	}
}

func TestUserCRUD(t *testing.T) {
	ctx := context.Background()
	db, cleanup := setupTestDB(ctx, t)
	defer cleanup()

	user := mockUser()

	t.Run("CreateUser", func(t *testing.T) {
		createUser(t, db, ctx, user)
	})

	t.Run("GetUserByEmail", func(t *testing.T) {
		getUserByEmail(t, db, ctx, user)
	})

	t.Run("GetUserByID", func(t *testing.T) {
		getUserByID(t, db, ctx, user)
	})

	t.Run("GetUserByUsername", func(t *testing.T) {
		getUserByUsername(t, db, ctx, user)
	})

	t.Run("UpdateUser", func(t *testing.T) {
		user = updateUser(t, db, ctx, user)
	})

	t.Run("DeleteUser", func(t *testing.T) {
		deleteUser(t, db, ctx, user.ID)
	})
}
