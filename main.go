package main

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/otakakot/sample-go-sqlc-sqldef/pkg/schema"
)

func main() {
	dsn := cmp.Or(os.Getenv("DATABSE_URL"), "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable")

	conn, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		panic(err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), conn)
	if err != nil {
		panic(err)
	}

	defer pool.Close()

	if err := pool.Ping(context.Background()); err != nil {
		panic(err)
	}

	ctx := context.Background()

	user, err := schema.New(pool).CreateUser(ctx, uuid.NewString())
	if err != nil {
		panic(err)
	}

	slog.InfoContext(ctx, fmt.Sprintf("user: %+v", user))

	user, err = schema.New(pool).UpdateUser(ctx, schema.UpdateUserParams{
		ID:   user.ID,
		Name: "updated",
	})
	if err != nil {
		panic(err)
	}

	slog.InfoContext(ctx, fmt.Sprintf("user: %+v", user))

	if err := schema.New(pool).DeleteUser(ctx, user.ID); err != nil {
		panic(err)
	}

	user, err = schema.New(pool).FindUserByID(ctx, user.ID)
	if err != nil {
		err = handle(err)
		if err != nil {
			panic(err)
		}
	}
}

func handle(err error) error {
	if errors.Is(pgx.ErrNoRows, err) {
		slog.Warn(err.Error())

		return nil
	}

	return err
}
