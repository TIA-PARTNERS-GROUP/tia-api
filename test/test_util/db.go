// test/test_util/db.go
package testutil

import (
	"fmt"
	"log"
	"testing"

	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var TestDB *gorm.DB

func SetupTestDB() {
	dsn := "root:tia-dev-password@tcp(127.0.0.1:3306)/tia-dev?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to test database: %v", err)
	}

	// Run migrations
	allModels := []interface{}{
		&models.User{}, &models.Business{}, &models.Project{}, &models.Skill{},
		&models.Publication{}, &models.Notification{}, &models.UserSkill{},
		&models.ProjectSkill{}, &models.ProjectMember{}, &models.BusinessConnection{},
		&models.BusinessTag{}, &models.UserSession{},
		&models.Feedback{}, &models.ProjectApplicant{}, &models.UserConfig{},
		&models.L2EResponse{}, &models.Subscription{}, &models.UserSubscription{},
		&models.UserDailyActivityProgress{}, &models.Event{}, &models.DailyActivity{},
		&models.DailyActivityEnrolment{}, &models.Region{}, &models.ProjectRegion{},
		&models.InferredConnection{},
	}
	if err := db.AutoMigrate(allModels...); err != nil {
		log.Fatalf("Failed to migrate database for tests: %v", err)
	}

	TestDB = db
}

func CleanupTestDB(t *testing.T, db *gorm.DB) {
	tables, err := db.Migrator().GetTables()
	assert.NoError(t, err)

	err = db.Exec("SET FOREIGN_KEY_CHECKS = 0;").Error
	assert.NoError(t, err)

	for _, table := range tables {
		err := db.Exec(fmt.Sprintf("TRUNCATE TABLE %s;", table)).Error
		assert.NoError(t, err)
	}

	err = db.Exec("SET FOREIGN_KEY_CHECKS = 1;").Error
	assert.NoError(t, err)
}
