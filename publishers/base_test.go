package publishers

import (
	"os"
	"testing"

	"github.com/xeviknal/background-commons/database"
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func setup() {
	setupTestDatabase()
}

func teardown() {
	destroyDatabase()
}

func setupTestDatabase() {
	database.SetConnectionConfig("test", "test", "test", "localhost")
}

func destroyDatabase() {
	database.Clean()
}
