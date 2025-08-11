package testx

import (
	"database/sql"
	"os"
	"path/filepath"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/go-testfixtures/testfixtures/v3"
)

func SetupPostgres(
	t *testing.T,
) string {
	t.Helper()

	user := "test"
	password := "test"
	db := "test"

	pwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	schemaPath := filepath.Join(filepath.Dir(pwd), "schema")

	if _, err := os.Stat(schemaPath); os.IsNotExist(err) {
		t.Fatalf("Schema directory does not exist: %s", schemaPath)
	}

	sqlFiles, err := filepath.Glob(filepath.Join(schemaPath, "*.sql"))
	if err != nil {
		t.Fatal(err)
	}

	if len(sqlFiles) == 0 {
		t.Fatalf("No SQL files found in schema directory: %s", schemaPath)
	}

	container, err := postgres.Run(
		t.Context(),
		"postgres:17-alpine",
		postgres.WithDatabase(db),
		postgres.WithUsername(user),
		postgres.WithPassword(password),
		postgres.WithInitScripts(sqlFiles...),
		testcontainers.WithEnv(map[string]string{
			"TZ":                        "UTC",
			"LANG":                      "ja_JP.UTF-8",
			"POSTGRES_INITDB_ARGS":      "--encoding=UTF-8",
			"POSTGRES_HOST_AUTH_METHOD": "trust",
		}),
		testcontainers.WithWaitStrategy(
			wait.ForAll(
				wait.ForListeningPort("5432/tcp"),
				wait.ForExec([]string{"pg_isready", "-U", user, "-d", db}).
					WithPollInterval(1*time.Second).
					WithExitCodeMatcher(func(exitCode int) bool {
						return exitCode == 0
					}).
					WithStartupTimeout(30*time.Second),
			),
		),
	)
	if err != nil {
		t.Fatal(err)
	}

	testcontainers.CleanupContainer(t, container)

	dsn, err := container.ConnectionString(t.Context(), "sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	return dsn
}

func SetupData(
	t *testing.T,
	dsn string,
) {
	t.Helper()

	dia := "postgres"

	db, err := sql.Open(dia, dsn)
	if err != nil {
		t.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		t.Fatal(err)
	}

	fixtures, err := testfixtures.New(
		testfixtures.Database(db),
		testfixtures.Dialect(dia),
		testfixtures.Directory("testdata"),
	)
	if err != nil {
		t.Fatal(err)
	}

	if err := fixtures.Load(); err != nil {
		t.Fatal(err)
	}
}
