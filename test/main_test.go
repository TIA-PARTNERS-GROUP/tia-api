// test/main_test.go
package main

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/models"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var testDB *gorm.DB

func TestMain(m *testing.M) {
	if err := godotenv.Load("../.env"); err != nil {
		log.Printf("Warning: could not load .env file. Relying on environment variables: %v", err)
	}

	dsn := "root:tia-dev-password@tcp(127.0.0.1:3306)/tia-dev?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to test database for TestMain: %v", err)
	}
	testDB = db

	allModels := []interface{}{
		&models.User{}, &models.Business{}, &models.Project{}, &models.Skill{},
		&models.Publication{}, &models.Idea{}, &models.Notification{}, &models.UserSkill{},
		&models.ProjectSkill{}, &models.ProjectMember{}, &models.BusinessConnection{},
		&models.BusinessTag{}, &models.IdeaVote{}, &models.UserSession{},
	}
	if err := testDB.AutoMigrate(allModels...); err != nil {
		log.Fatalf("Failed to migrate database for TestMain: %v", err)
	}

	exitCode := m.Run()

	os.Exit(exitCode)
}

func cleanupTestDB(t *testing.T) {
	tables, err := testDB.Migrator().GetTables()
	assert.NoError(t, err)

	err = testDB.Exec("SET FOREIGN_KEY_CHECKS = 0;").Error
	assert.NoError(t, err)

	for _, table := range tables {
		err := testDB.Exec(fmt.Sprintf("TRUNCATE TABLE %s;", table)).Error
		assert.NoError(t, err)
	}

	err = testDB.Exec("SET FOREIGN_KEY_CHECKS = 1;").Error
	assert.NoError(t, err)
}
