package test_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/otakakot/sample-go-sqlc-sqldef/pkg/schema"
	"github.com/otakakot/sample-go-sqlc-sqldef/test/testx"
)

func TestCreateUser(t *testing.T) {
	t.Parallel()

	dsn, cleanup, err := testx.SetupDB(t)
	if err != nil {
		t.Fatalf("failed to setup db: %v", err)
	}

	t.Cleanup(cleanup)

	conn, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		t.Fatal()
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), conn)
	if err != nil {
		t.Fatal()
	}

	t.Cleanup(pool.Close)

	if err := pool.Ping(context.Background()); err != nil {
		t.Fatalf("failed to ping db: %v", err)
	}

	ctx := context.Background()

	name := uuid.NewString()

	user, err := schema.New(pool).CreateUser(ctx, name)
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	if user.Name != name {
		t.Errorf("user name is not correct: %s", user.Name)
	}
}

func TestFindUserByID(t *testing.T) {
	t.Parallel()

	dsn, cleanup, err := testx.SetupDB(t)
	if err != nil {
		t.Fatalf("failed to setup db: %v", err)
	}

	t.Cleanup(cleanup)

	conn, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		t.Fatal()
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), conn)
	if err != nil {
		t.Fatal()
	}

	t.Cleanup(pool.Close)

	if err := testx.SetupData(t, dsn); err != nil {
		t.Fatalf("failed to setup data: %v", err)
	}

	id := uuid.MustParse("77777777-7777-7777-7777-777777777777")

	user, err := schema.New(pool).FindUserByID(context.Background(), id)
	if err != nil {
		t.Fatalf("failed to find user: %v", err)
	}

	if user.ID != id {
		t.Errorf("user id is not correct: %s", user.ID)
	}

	if _, err := schema.New(pool).FindUserByID(context.Background(), uuid.New()); err == nil {
		t.Errorf("user should not be found")
	} else {
		if !errors.Is(pgx.ErrNoRows, err) {
			t.Errorf("unexpected error: %v", err)
		}
	}
}
