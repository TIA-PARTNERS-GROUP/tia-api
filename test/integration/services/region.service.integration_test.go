package main

import (
	"context"
	"testing"

	"github.com/TIA-PARTNERS-GROUP/tia-api/internal/core/services"
	testutil "github.com/TIA-PARTNERS-GROUP/tia-api/test/test_util"
	"github.com/stretchr/testify/assert"
)

func TestRegionService_Integration_SeedAndGet(t *testing.T) {
	testutil.CleanupTestDB(t, testutil.TestDB)
	regionService := services.NewRegionService(testutil.TestDB)

	err := regionService.SeedRegions(context.Background())
	assert.NoError(t, err)

	regions, err := regionService.GetAllRegions(context.Background())
	assert.NoError(t, err)
	assert.NotEmpty(t, regions)
	assert.Len(t, regions, 8)

	assert.Equal(t, "Australian Capital Territory", regions[0].Name)

	err = regionService.SeedRegions(context.Background())
	assert.NoError(t, err)
	regions, _ = regionService.GetAllRegions(context.Background())
	assert.Len(t, regions, 8)
}
