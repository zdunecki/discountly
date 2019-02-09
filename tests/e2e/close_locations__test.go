package e2e

import (
	"github.com/stretchr/testify/assert"
	"github.com/zdunecki/discountly/features/discounts/models"
	"github.com/zdunecki/discountly/features/search/finder"
	"github.com/zdunecki/discountly/infra"
	"github.com/zdunecki/discountly/tests"
	"testing"
)

var dbAddress string
var geoAddress string

func init() {
	dbAddress = infra.GetEnv("DB_CONNECTION")
	geoAddress = infra.GetEnv("GEO_CONNECTION")
}

func TestCloseLocations(t *testing.T) {
	assert := assert.New(t)

	defer tests.DeleteAll()

	l := discounts.Location{
		Lat: testLat,
		Lon: testLon,
	}.New()

	discountLocations := []discounts.Location{
		l,
	}

	_ = finder.SetLocationPoint(discountLocations[0].Id, discountLocations)

	result := finder.CloseLocations(discountLocations[0].Id, discounts.Location{
		Lat: nearbyTestLat,
		Lon: nearbyTestLon,
	}, discountLocations)

	assert.Equal(result, []discounts.Location{
		{
			Id:  discountLocations[0].Id + "|" + l.Id,
			Lat: testLat,
			Lon: testLon,
		},
	})
}
