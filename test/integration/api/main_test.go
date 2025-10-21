package main
import (
	"os"
	"testing"
	testutil "github.com/TIA-PARTNERS-GROUP/tia-api/test/test_util"
	"github.com/joho/godotenv"
)
func TestMain(m *testing.M) {
	
	godotenv.Load("../../../.env")
	
	testutil.SetupTestDB()
	
	exitCode := m.Run()
	os.Exit(exitCode)
}
