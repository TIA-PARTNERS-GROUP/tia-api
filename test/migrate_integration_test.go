package main

import (
	"testing"

	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/models"
)

func TestIntegration_SchemaCreation(t *testing.T) {
	if !testDB.Migrator().HasTable(&models.User{}) {
		t.Errorf("Table 'users' was not created by TestMain migration")
	}
	if !testDB.Migrator().HasTable("projects") {
		t.Errorf("Table 'projects' was not created by TestMain migration")
	}
	if !testDB.Migrator().HasTable("business_connections") {
		t.Errorf("Table 'business_connections' was not created by TestMain migration")
	}

	t.Log("Integration test: Schema creation verified successfully.")
}
