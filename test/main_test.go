package test_test

import (
	"os"
	"testing"

	"github.com/otakakot/sample-go-sqlc-sqldef/test/testx"
)

func TestMain(m *testing.M) {
	_, cleanup, err := testx.GlovalDSN(&testing.T{})
	if err != nil {
		panic(err)
	}

	code := m.Run()

	cleanup()

	os.Exit(code)
}
