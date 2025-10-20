package main

import (
	"os"
	"testing"

	testutil "github.com/TIA-PARTNERS-GROUP/tia-api/test/test_util"
	"github.com/joho/godotenv"
)

// TestMain runs once for the entire 'api' test package.
func TestMain(m *testing.M) {
	// Go up two directories to find the .env file from the project root.
	godotenv.Load("../../../.env")

	// Call the shared setup function. This initializes test_util.TestDB.
	testutil.SetupTestDB()

	// Run all tests in this package.
	exitCode := m.Run()

	os.Exit(exitCode)
}
