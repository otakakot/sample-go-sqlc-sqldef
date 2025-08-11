package test_test

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/otakakot/sample-go-sqlc-sqldef/pkg/schema"
	"github.com/otakakot/sample-go-sqlc-sqldef/test/internal/testx"
)

func TestCreateUser(t *testing.T) {
	t.Parallel()

	dsn := testx.SetupPostgres(t)

	conn, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		t.Fatalf("failed to parse config: %v", err)
	}

	pool, err := pgxpool.NewWithConfig(t.Context(), conn)
	if err != nil {
		t.Fatalf("failed to create pool: %v", err)
	}

	t.Cleanup(pool.Close)

	if err := pool.Ping(t.Context()); err != nil {
		t.Fatalf("failed to ping db: %v", err)
	}

	name := uuid.NewString()

	user, err := schema.New(pool).CreateUser(t.Context(), name)
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	if user.Name != name {
		t.Errorf("user name is not correct: %s", user.Name)
	}
}

func TestFindUserByID(t *testing.T) {
	t.Parallel()

	dsn := testx.SetupPostgres(t)

	conn, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		t.Fatalf("failed to parse config: %v", err)
	}

	pool, err := pgxpool.NewWithConfig(t.Context(), conn)
	if err != nil {
		t.Fatalf("failed to create pool: %v", err)
	}

	t.Cleanup(pool.Close)

	testx.SetupData(t, dsn)

	id := uuid.MustParse("77777777-7777-7777-7777-777777777777")

	user, err := schema.New(pool).FindUserByID(t.Context(), id)
	if err != nil {
		t.Fatalf("failed to find user: %v", err)
	}

	if user.ID != id {
		t.Errorf("user id is not correct: %s", user.ID)
	}

	if _, err := schema.New(pool).FindUserByID(t.Context(), uuid.New()); err == nil {
		t.Errorf("user should not be found")
	} else {
		if !errors.Is(pgx.ErrNoRows, err) {
			t.Errorf("unexpected error: %v", err)
		}
	}
}

func TestUpdateUser(t *testing.T) {
	t.Parallel()

	dsn := testx.SetupPostgres(t)

	conn, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		t.Fatalf("failed to parse config: %v", err)
	}

	pool, err := pgxpool.NewWithConfig(t.Context(), conn)
	if err != nil {
		t.Fatalf("failed to create pool: %v", err)
	}

	t.Cleanup(pool.Close)

	testx.SetupData(t, dsn)

	id := uuid.MustParse("77777777-7777-7777-7777-777777777777")

	want, err := schema.New(pool).FindUserByID(t.Context(), id)
	if err != nil {
		t.Fatalf("failed to find user: %v", err)
	}

	name := uuid.NewString()

	got, err := schema.New(pool).UpdateUser(t.Context(), schema.UpdateUserParams{
		ID:   id,
		Name: name,
	})
	if err != nil {
		t.Fatalf("failed to update user: %v", err)
	}

	if got.ID != want.ID {
		t.Errorf("user id is not correct: %s", got.ID)
	}

	if got.Name != name {
		t.Errorf("user name is not correct: %s", got.Name)
	}

	if got.CreatedAt != want.CreatedAt {
		t.Errorf("user created_at is not correct: %s", got.CreatedAt.Time)
	}

	if got.UpdatedAt == want.UpdatedAt {
		t.Errorf("user updated_at is not updated")
	}
}

func TestDeleteUser(t *testing.T) {
	t.Parallel()

	dsn := testx.SetupPostgres(t)

	conn, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		t.Fatalf("failed to parse config: %v", err)
	}

	pool, err := pgxpool.NewWithConfig(t.Context(), conn)
	if err != nil {
		t.Fatalf("failed to create pool: %v", err)
	}

	t.Cleanup(pool.Close)

	testx.SetupData(t, dsn)

	id := uuid.MustParse("77777777-7777-7777-7777-777777777777")

	if err := schema.New(pool).DeleteUser(t.Context(), id); err != nil {
		t.Fatalf("failed to delete user: %v", err)
	}

	if _, err := schema.New(pool).FindUserByID(t.Context(), id); err == nil {
		t.Errorf("user should not be found")
	} else {
		if !errors.Is(pgx.ErrNoRows, err) {
			t.Errorf("unexpected error: %v", err)
		}
	}
}
