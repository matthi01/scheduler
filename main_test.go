package main

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	loadEnv()
	a.Initialize(
		os.Getenv("HOST"),
		os.Getenv("APP_DB_PORT"),
		os.Getenv("APP_DB_USERNAME"),
		os.Getenv("APP_TEST_DB_NAME"),
	)

	ensureTablesExists()
	code := m.Run()
	clearTables()
	os.Exit(code)
}
