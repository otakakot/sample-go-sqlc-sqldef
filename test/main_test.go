package test_test

import (
	"context"
	"log/slog"
	"os"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		panic(err)
	}

	if err := pool.Client.Ping(); err != nil {
		panic(err)
	}

	// test ファイルが置かれているディレクトリ
	pwd, _ := os.Getwd()

	ddl := strings.Replace(pwd, "test", "schema", 1)

	slog.Info(ddl)

	opt := dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "16-alpine",
		Env: []string{
			"POSTGRES_USER=postgres",
			"POSTGRES_PASSWORD=postgres",
			"POSTGRES_DB=postgres",
			"listen_addresses='*'",
		},
		Mounts: []string{
			ddl + ":/docker-entrypoint-initdb.d",
		},
	}

	hcOpt := func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	}

	hcOpts := []func(*docker.HostConfig){
		hcOpt,
	}

	resource, err := pool.RunWithOptions(&opt, hcOpts...)
	if err != nil {
		panic(err)
	}

	port := resource.GetPort("5432/tcp")

	dsn := "postgres://postgres:postgres@localhost:" + port + "/postgres?sslmode=disable"

	slog.Info(dsn)

	if err := pool.Retry(func() error {
		conn, err := pgxpool.ParseConfig(dsn)
		if err != nil {
			return err
		}

		pool, err := pgxpool.NewWithConfig(context.Background(), conn)
		if err != nil {
			return err
		}

		if err := pool.Ping(context.Background()); err != nil {
			return err
		}

		return nil
	}); err != nil {
		panic(err)
	}

	code := m.Run()

	if err := pool.Purge(resource); err != nil {
		panic(err)
	}

	os.Exit(code)
}
