package main
import (
	"context"
	"testing"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/core/services"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/models"
	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/ports"
	testutil "github.com/TIA-PARTNERS-GROUP/tia-api/test/test_util"
	"github.com/stretchr/testify/assert"
)
func TestBusinessConnectionService_Integration_CreateAndGet(t *testing.T) {
	testutil.CleanupTestDB(t, testutil.TestDB)
	businessConnectionService := services.NewBusinessConnectionService(testutil.TestDB)
	user := models.User{FirstName: "Connector", LoginEmail: "connector@business.com", Active: true}
	result := testutil.TestDB.Create(&user)
	assert.NoError(t, result.Error)
	business1 := models.Business{
		Name:             "Business 1",
		OperatorUserID:   user.ID,
		BusinessType:     models.BusinessTypeTechnology,
		BusinessCategory: models.BusinessCategoryB2B,
		BusinessPhase:    models.BusinessPhaseStartup,
		Active:           true,
	}
	result = testutil.TestDB.Create(&business1)
	assert.NoError(t, result.Error)
	business2 := models.Business{
		Name:             "Business 2",
		OperatorUserID:   user.ID,
		BusinessType:     models.BusinessTypeConsulting,
		BusinessCategory: models.BusinessCategoryB2B,
		BusinessPhase:    models.BusinessPhaseStartup,
		Active:           true,
	}
	result = testutil.TestDB.Create(&business2)
	assert.NoError(t, result.Error)
	notes := "Looking to partner on new projects"
	createDTO := ports.CreateBusinessConnectionInput{
		InitiatingBusinessID: business1.ID,
		ReceivingBusinessID:  business2.ID,
		ConnectionType:       models.ConnectionTypePartnership,
		InitiatedByUserID:    user.ID,
		Notes:                &notes,
	}
	createdConnection, err := businessConnectionService.CreateBusinessConnection(context.Background(), createDTO)
	assert.NoError(t, err)
	assert.NotNil(t, createdConnection)
	assert.Equal(t, business1.ID, createdConnection.InitiatingBusinessID)
	assert.Equal(t, business2.ID, createdConnection.ReceivingBusinessID)
	assert.Equal(t, models.ConnectionTypePartnership, createdConnection.ConnectionType)
	assert.Equal(t, models.ConnectionStatusPending, createdConnection.Status)
	assert.Equal(t, user.ID, createdConnection.InitiatedByUserID)
	fetchedConnection, err := businessConnectionService.GetBusinessConnection(context.Background(), createdConnection.ID)
	assert.NoError(t, err)
	assert.NotNil(t, fetchedConnection)
	assert.Equal(t, createdConnection.ConnectionType, fetchedConnection.ConnectionType)
	assert.Equal(t, createdConnection.Status, fetchedConnection.Status)
}
func TestBusinessConnectionService_Integration_DuplicatePrevention(t *testing.T) {
	testutil.CleanupTestDB(t, testutil.TestDB)
	businessConnectionService := services.NewBusinessConnectionService(testutil.TestDB)
	user := models.User{FirstName: "Connector", LoginEmail: "connector2@business.com", Active: true}
	testutil.TestDB.Create(&user)
	business1 := models.Business{
		Name:             "Business A",
		OperatorUserID:   user.ID,
		BusinessType:     models.BusinessTypeRetail,
		BusinessCategory: models.BusinessCategoryB2B,
		BusinessPhase:    models.BusinessPhaseStartup,
		Active:           true,
	}
	testutil.TestDB.Create(&business1)
	business2 := models.Business{
		Name:             "Business B",
		OperatorUserID:   user.ID,
		BusinessType:     models.BusinessTypeServices,
		BusinessCategory: models.BusinessCategoryB2B,
		BusinessPhase:    models.BusinessPhaseStartup,
		Active:           true,
	}
	testutil.TestDB.Create(&business2)
	createDTO := ports.CreateBusinessConnectionInput{
		InitiatingBusinessID: business1.ID,
		ReceivingBusinessID:  business2.ID,
		ConnectionType:       models.ConnectionTypeSupplier,
		InitiatedByUserID:    user.ID,
	}
	_, err := businessConnectionService.CreateBusinessConnection(context.Background(), createDTO)
	assert.NoError(t, err)
	_, err = businessConnectionService.CreateBusinessConnection(context.Background(), createDTO)
	assert.Error(t, err)
	assert.Equal(t, ports.ErrBusinessConnectionAlreadyExists, err)
}
func TestBusinessConnectionService_Integration_CannotConnectToSelf(t *testing.T) {
	testutil.CleanupTestDB(t, testutil.TestDB)
	businessConnectionService := services.NewBusinessConnectionService(testutil.TestDB)
	user := models.User{FirstName: "Connector", LoginEmail: "connector3@business.com", Active: true}
	testutil.TestDB.Create(&user)
	business := models.Business{
		Name:             "Business Self",
		OperatorUserID:   user.ID,
		BusinessType:     models.BusinessTypeTechnology,
		BusinessCategory: models.BusinessCategoryB2B,
		BusinessPhase:    models.BusinessPhaseStartup,
		Active:           true,
	}
	testutil.TestDB.Create(&business)
	createDTO := ports.CreateBusinessConnectionInput{
		InitiatingBusinessID: business.ID,
		ReceivingBusinessID:  business.ID,
		ConnectionType:       models.ConnectionTypePartnership,
		InitiatedByUserID:    user.ID,
	}
	_, err := businessConnectionService.CreateBusinessConnection(context.Background(), createDTO)
	assert.Error(t, err)
	assert.Equal(t, ports.ErrCannotConnectToSelf, err)
}
func TestBusinessConnectionService_Integration_UpdateAndDelete(t *testing.T) {
	testutil.CleanupTestDB(t, testutil.TestDB)
	businessConnectionService := services.NewBusinessConnectionService(testutil.TestDB)
	user := models.User{FirstName: "Connector", LoginEmail: "connector4@business.com", Active: true}
	testutil.TestDB.Create(&user)
	business1 := models.Business{
		Name:             "Business X",
		OperatorUserID:   user.ID,
		BusinessType:     models.BusinessTypeManufacturing,
		BusinessCategory: models.BusinessCategoryB2B,
		BusinessPhase:    models.BusinessPhaseStartup,
		Active:           true,
	}
	testutil.TestDB.Create(&business1)
	business2 := models.Business{
		Name:             "Business Y",
		OperatorUserID:   user.ID,
		BusinessType:     models.BusinessTypeConsulting,
		BusinessCategory: models.BusinessCategoryB2B,
		BusinessPhase:    models.BusinessPhaseStartup,
		Active:           true,
	}
	testutil.TestDB.Create(&business2)
	createDTO := ports.CreateBusinessConnectionInput{
		InitiatingBusinessID: business1.ID,
		ReceivingBusinessID:  business2.ID,
		ConnectionType:       models.ConnectionTypeClient,
		InitiatedByUserID:    user.ID,
	}
	connection, err := businessConnectionService.CreateBusinessConnection(context.Background(), createDTO)
	assert.NoError(t, err)
	newNotes := "Updated notes about this connection"
	collaborationType := models.ConnectionTypeCollaboration
	updateDTO := ports.UpdateBusinessConnectionInput{
		ConnectionType: &collaborationType,
		Notes:          &newNotes,
	}
	updatedConnection, err := businessConnectionService.UpdateBusinessConnection(context.Background(), connection.ID, updateDTO)
	assert.NoError(t, err)
	assert.NotNil(t, updatedConnection)
	assert.Equal(t, models.ConnectionTypeCollaboration, updatedConnection.ConnectionType)
	assert.Equal(t, newNotes, *updatedConnection.Notes)
	err = businessConnectionService.DeleteBusinessConnection(context.Background(), connection.ID)
	assert.NoError(t, err)
	_, err = businessConnectionService.GetBusinessConnection(context.Background(), connection.ID)
	assert.Error(t, err)
	assert.Equal(t, ports.ErrBusinessConnectionNotFound, err)
}
func TestBusinessConnectionService_Integration_AcceptAndReject(t *testing.T) {
	testutil.CleanupTestDB(t, testutil.TestDB)
	businessConnectionService := services.NewBusinessConnectionService(testutil.TestDB)
	user := models.User{FirstName: "Connector", LoginEmail: "connector5@business.com", Active: true}
	testutil.TestDB.Create(&user)
	business1 := models.Business{
		Name:             "Business Alpha",
		OperatorUserID:   user.ID,
		BusinessType:     models.BusinessTypeTechnology,
		BusinessCategory: models.BusinessCategoryB2B,
		BusinessPhase:    models.BusinessPhaseStartup,
		Active:           true,
	}
	testutil.TestDB.Create(&business1)
	business2 := models.Business{
		Name:             "Business Beta",
		OperatorUserID:   user.ID,
		BusinessType:     models.BusinessTypeConsulting,
		BusinessCategory: models.BusinessCategoryB2B,
		BusinessPhase:    models.BusinessPhaseStartup,
		Active:           true,
	}
	testutil.TestDB.Create(&business2)
	createDTO := ports.CreateBusinessConnectionInput{
		InitiatingBusinessID: business1.ID,
		ReceivingBusinessID:  business2.ID,
		ConnectionType:       models.ConnectionTypePartnership,
		InitiatedByUserID:    user.ID,
	}
	connection, err := businessConnectionService.CreateBusinessConnection(context.Background(), createDTO)
	assert.NoError(t, err)
	assert.Equal(t, models.ConnectionStatusPending, connection.Status)
	acceptedConnection, err := businessConnectionService.AcceptBusinessConnection(context.Background(), connection.ID)
	assert.NoError(t, err)
	assert.NotNil(t, acceptedConnection)
	assert.Equal(t, models.ConnectionStatusActive, acceptedConnection.Status)
	createDTO2 := ports.CreateBusinessConnectionInput{
		InitiatingBusinessID: business2.ID,
		ReceivingBusinessID:  business1.ID,
		ConnectionType:       models.ConnectionTypeReferral,
		InitiatedByUserID:    user.ID,
	}
	connection2, err := businessConnectionService.CreateBusinessConnection(context.Background(), createDTO2)
	assert.NoError(t, err)
	rejectedConnection, err := businessConnectionService.RejectBusinessConnection(context.Background(), connection2.ID)
	assert.NoError(t, err)
	assert.NotNil(t, rejectedConnection)
	assert.Equal(t, models.ConnectionStatusRejected, rejectedConnection.Status)
}
func TestBusinessConnectionService_Integration_GetConnections(t *testing.T) {
	testutil.CleanupTestDB(t, testutil.TestDB)
	businessConnectionService := services.NewBusinessConnectionService(testutil.TestDB)
	user := models.User{FirstName: "Connector", LoginEmail: "connector6@business.com", Active: true}
	testutil.TestDB.Create(&user)
	business1 := models.Business{
		Name:             "Business One",
		OperatorUserID:   user.ID,
		BusinessType:     models.BusinessTypeTechnology,
		BusinessCategory: models.BusinessCategoryB2B,
		BusinessPhase:    models.BusinessPhaseStartup,
		Active:           true,
	}
	testutil.TestDB.Create(&business1)
	business2 := models.Business{
		Name:             "Business Two",
		OperatorUserID:   user.ID,
		BusinessType:     models.BusinessTypeConsulting,
		BusinessCategory: models.BusinessCategoryB2B,
		BusinessPhase:    models.BusinessPhaseStartup,
		Active:           true,
	}
	testutil.TestDB.Create(&business2)
	business3 := models.Business{
		Name:             "Business Three",
		OperatorUserID:   user.ID,
		BusinessType:     models.BusinessTypeRetail,
		BusinessCategory: models.BusinessCategoryB2B,
		BusinessPhase:    models.BusinessPhaseStartup,
		Active:           true,
	}
	testutil.TestDB.Create(&business3)
	_, err := businessConnectionService.CreateBusinessConnection(context.Background(), ports.CreateBusinessConnectionInput{
		InitiatingBusinessID: business1.ID,
		ReceivingBusinessID:  business2.ID,
		ConnectionType:       models.ConnectionTypePartnership,
		InitiatedByUserID:    user.ID,
	})
	assert.NoError(t, err)
	_, err = businessConnectionService.CreateBusinessConnection(context.Background(), ports.CreateBusinessConnectionInput{
		InitiatingBusinessID: business1.ID,
		ReceivingBusinessID:  business3.ID,
		ConnectionType:       models.ConnectionTypeSupplier,
		InitiatedByUserID:    user.ID,
	})
	assert.NoError(t, err)
	_, err = businessConnectionService.CreateBusinessConnection(context.Background(), ports.CreateBusinessConnectionInput{
		InitiatingBusinessID: business2.ID,
		ReceivingBusinessID:  business1.ID,
		ConnectionType:       models.ConnectionTypeClient,
		InitiatedByUserID:    user.ID,
	})
	assert.NoError(t, err)
	connections, err := businessConnectionService.GetBusinessConnections(context.Background(), business1.ID, nil, nil)
	assert.NoError(t, err)
	assert.Len(t, connections, 3)
	pendingConnections, err := businessConnectionService.GetPendingConnections(context.Background(), business2.ID)
	assert.NoError(t, err)
	assert.Len(t, pendingConnections, 1)
	activeStatus := models.ConnectionStatusActive
	activeConnections, err := businessConnectionService.GetBusinessConnections(context.Background(), business1.ID, nil, &activeStatus)
	assert.NoError(t, err)
	assert.Len(t, activeConnections, 0)
}
func TestBusinessConnectionService_Integration_Validation(t *testing.T) {
	testutil.CleanupTestDB(t, testutil.TestDB)
	businessConnectionService := services.NewBusinessConnectionService(testutil.TestDB)
	createDTO := ports.CreateBusinessConnectionInput{
		InitiatingBusinessID: 999,
		ReceivingBusinessID:  1,
		ConnectionType:       models.ConnectionTypePartnership,
		InitiatedByUserID:    1,
	}
	_, err := businessConnectionService.CreateBusinessConnection(context.Background(), createDTO)
	assert.Error(t, err)
	assert.Equal(t, ports.ErrBusinessNotFound, err)
	user := models.User{FirstName: "Connector", LoginEmail: "connector7@business.com", Active: true}
	testutil.TestDB.Create(&user)
	business1 := models.Business{
		Name:             "Business Test",
		OperatorUserID:   user.ID,
		BusinessType:     models.BusinessTypeTechnology,
		BusinessCategory: models.BusinessCategoryB2B,
		BusinessPhase:    models.BusinessPhaseStartup,
		Active:           true,
	}
	testutil.TestDB.Create(&business1)
	business2 := models.Business{
		Name:             "Business Test 2",
		OperatorUserID:   user.ID,
		BusinessType:     models.BusinessTypeConsulting,
		BusinessCategory: models.BusinessCategoryB2B,
		BusinessPhase:    models.BusinessPhaseStartup,
		Active:           true,
	}
	testutil.TestDB.Create(&business2)
	createDTO = ports.CreateBusinessConnectionInput{
		InitiatingBusinessID: business1.ID,
		ReceivingBusinessID:  business2.ID,
		ConnectionType:       models.ConnectionTypePartnership,
		InitiatedByUserID:    999,
	}
	_, err = businessConnectionService.CreateBusinessConnection(context.Background(), createDTO)
	assert.Error(t, err)
	assert.Equal(t, ports.ErrUserNotFound, err)
}
func TestBusinessConnectionService_Integration_NonExistentConnection(t *testing.T) {
	testutil.CleanupTestDB(t, testutil.TestDB)
	businessConnectionService := services.NewBusinessConnectionService(testutil.TestDB)
	_, err := businessConnectionService.GetBusinessConnection(context.Background(), 999)
	assert.Error(t, err)
	assert.Equal(t, ports.ErrBusinessConnectionNotFound, err)
	partnershipType := models.ConnectionTypePartnership
	updateDTO := ports.UpdateBusinessConnectionInput{
		ConnectionType: &partnershipType,
	}
	_, err = businessConnectionService.UpdateBusinessConnection(context.Background(), 999, updateDTO)
	assert.Error(t, err)
	assert.Equal(t, ports.ErrBusinessConnectionNotFound, err)
	err = businessConnectionService.DeleteBusinessConnection(context.Background(), 999)
	assert.Error(t, err)
	assert.Equal(t, ports.ErrBusinessConnectionNotFound, err)
}
func TestBusinessConnectionService_Integration_AcceptNonPendingConnection(t *testing.T) {
	testutil.CleanupTestDB(t, testutil.TestDB)
	businessConnectionService := services.NewBusinessConnectionService(testutil.TestDB)
	user := models.User{FirstName: "Connector", LoginEmail: "connector8@business.com", Active: true}
	testutil.TestDB.Create(&user)
	business1 := models.Business{
		Name:             "Business A",
		OperatorUserID:   user.ID,
		BusinessType:     models.BusinessTypeTechnology,
		BusinessCategory: models.BusinessCategoryB2B,
		BusinessPhase:    models.BusinessPhaseStartup,
		Active:           true,
	}
	testutil.TestDB.Create(&business1)
	business2 := models.Business{
		Name:             "Business B",
		OperatorUserID:   user.ID,
		BusinessType:     models.BusinessTypeConsulting,
		BusinessCategory: models.BusinessCategoryB2B,
		BusinessPhase:    models.BusinessPhaseStartup,
		Active:           true,
	}
	testutil.TestDB.Create(&business2)
	createDTO := ports.CreateBusinessConnectionInput{
		InitiatingBusinessID: business1.ID,
		ReceivingBusinessID:  business2.ID,
		ConnectionType:       models.ConnectionTypePartnership,
		InitiatedByUserID:    user.ID,
	}
	connection, err := businessConnectionService.CreateBusinessConnection(context.Background(), createDTO)
	assert.NoError(t, err)
	_, err = businessConnectionService.AcceptBusinessConnection(context.Background(), connection.ID)
	assert.NoError(t, err)
	_, err = businessConnectionService.AcceptBusinessConnection(context.Background(), connection.ID)
	assert.Error(t, err)
	assert.Equal(t, ports.ErrConnectionNotPending, err)
}
