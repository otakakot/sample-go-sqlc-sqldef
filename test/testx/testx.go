package testx

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	_ "github.com/lib/pq"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

func SetupDB(
	t *testing.T,
) (string, func(), error) {
	t.Helper()

	pool, err := dockertest.NewPool("")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create new pool: %w", err)
	}

	if err := pool.Client.Ping(); err != nil {
		return "", nil, fmt.Errorf("failed to ping: %w", err)
	}

	pwd, _ := os.Getwd()

	ddl := strings.Replace(pwd, "test", "schema", 1)

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
		return "", nil, fmt.Errorf("failed to run with options: %w", err)
	}

	port := resource.GetPort("5432/tcp")

	dsn := "postgres://postgres:postgres@localhost:" + port + "/postgres?sslmode=disable"

	if err := pool.Retry(func() error {
		conn, err := pgxpool.ParseConfig(dsn)
		if err != nil {
			return fmt.Errorf("failed to parse config: %w", err)
		}

		pool, err := pgxpool.NewWithConfig(context.Background(), conn)
		if err != nil {
			return fmt.Errorf("failed to create pool: %w", err)
		}

		if err := pool.Ping(context.Background()); err != nil {
			return fmt.Errorf("failed to ping: %w", err)
		}

		return nil
	}); err != nil {
		return "", nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}

	return dsn, func() {
		if err := pool.Purge(resource); err != nil {
			t.Log("failed to purge resource. error: " + err.Error())
		}
	}, nil
}
