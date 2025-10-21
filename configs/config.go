package configs
import (
	"log"
	"os"
	"github.com/joho/godotenv"
)
type Config struct {
	DatabaseURL string
}
func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}
	return &Config{
		DatabaseURL: dbURL,
	}
}
