package e2e

import (
	"github.com/zdunecki/discountly/features/auth/models"
	"github.com/zdunecki/discountly/features/discounts/models"
	"github.com/zdunecki/discountly/features/search/finder"
	"github.com/zdunecki/discountly/features/search/models"
	"github.com/zdunecki/discountly/infra"
	"github.com/zdunecki/discountly/tests"
	"github.com/stretchr/testify/assert"
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

	locations := []discounts.Location{
		l,
	}

	user := auth.User{}.New()
	finder.SetPoint(user.Id, locations)
	result := finder.CloseLocations(search.Search{
		Location: discounts.Location{
			Lat: nearbyTestLat,
			Lon: nearbyTestLon,
		},
	}, locations)
	assert.Equal(result, []discounts.Location{
		{
			Id:  l.Id,
			Lat: testLat,
			Lon: testLon,
		},
	})
}
