package finder

import (
	"github.com/zdunecki/discountly/features/discounts/models"
	"github.com/zdunecki/discountly/features/search/models"
	"github.com/zdunecki/discountly/infra"
	"github.com/gomodule/redigo/redis"
	"github.com/paulmach/go.geojson"
)

const km = 1000
const nearbyLocation = 3 * km

var address string

func init() {
	address = infra.GetEnv("GEO_CONNECTION")
}

func SetPoint(userId string, locations []discounts.Location) {
	conn, err := redis.Dial("tcp", address)
	if err != nil {
		panic(err)
	}

	defer conn.Close()

	for _, l := range locations {
		g := geojson.NewPointGeometry([]float64{l.Lon, l.Lat})
		g.Type = "Point"

		rawJSON, err := g.MarshalJSON()
		if err != nil {
			panic(err)
		}

		_, err = conn.Do("SET", "location", l.Id, "OBJECT", rawJSON)
		if err != nil {
			panic(err)
		}
	}
}

func CloseLocations(search search.Search, discountLocations []discounts.Location) []discounts.Location {
	var locations []discounts.Location
	conn, err := redis.Dial("tcp", address)
	if err != nil {
		panic(err)
	}

	defer conn.Close()

	values, err := redis.Values(
		conn.Do("NEARBY", "location", "POINT", search.Location.Lat, search.Location.Lon, nearbyLocation),
	)
	if err != nil {
		panic(err)
	}

	if len(values) < 2 {
		panic("invalid value")
	}

	// 0 - cursor
	values, err = redis.Values(values[1], nil)
	if err != nil {
		panic("invalid value")
	}

	for _, val := range values {
		strings, err := redis.Strings(val, nil)
		if err != nil || len(strings) < 2 {
			panic(err)
		}

		geometry, err := geojson.UnmarshalGeometry([]byte(strings[1]))
		if err != nil || len(strings) < 2 {
			panic(err)
		}
		locationId := strings[0]
		point := geometry.Point

		for _, discountLocation := range discountLocations {
			if discountLocation.Id != locationId {
				continue
			}

			locations = append(locations, discounts.Location{
				Id:  locationId,
				Lon: point[0],
				Lat: point[1],
			})
		}
	}

	return locations
}

//TODO: implement hooks
//func SetLocationHooks(userId string, locations []discounts.Location) {
//	conn, err := redis.Dial("tcp", ":9851")
//	if err != nil {
//		panic(err)
//	}
//
//	defer conn.Close()
//	appURL := infra.GetEnv("APP_URL")
//
//	for _, l := range locations {
//		keyId := userId + ":" + l.Id
//
//		endpoint := appURL + "/test-geo"
//		_, err = conn.Do("SETHOOK", keyId, endpoint, "NEARBY", "fleet", "FENCE", "POINT", l.Lat, l.Lon)
//		if err != nil {
//			panic(err)
//		}
//	}
//}
