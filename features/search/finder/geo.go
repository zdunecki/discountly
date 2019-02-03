package finder

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/paulmach/go.geojson"
	"github.com/zdunecki/discountly/db"
	"github.com/zdunecki/discountly/features/discounts/models"
	"github.com/zdunecki/discountly/infra"
)

const PointKey = "location"
const HookPointKey = "location.hook"
const GeoFenceUsersKey = "location.users"

const km = 1000
const nearbyLocation = 3 * km

func NewGeoFenceUserId(geoFenceId string) string {
	return GeoFenceUsersKey + "|" + geoFenceId
}

func newGeoFenceId(discountId string, location discounts.Location) string {
	lon := fmt.Sprintf("%f", location.Lon)
	lat := fmt.Sprintf("%f", location.Lat)

	return discountId + "|" + lat + "_" + lon
}

func newLocationId(discountId string, location discounts.Location) string {
	return discountId + "|" + location.Id
}

func setPoint(keyId, keyValue string, location discounts.Location) error {
	conn := database.GeoRedisPool.Get()

	defer func() {
		err := conn.Close()
		if err != nil {
			panic(err)
		}
	}()

	g := geojson.NewPointGeometry([]float64{location.Lon, location.Lat})
	g.Type = "Point"

	rawJSON, err := g.MarshalJSON()
	if err != nil {
		return err
	}

	_, err = conn.Do("SET", keyId, keyValue, "OBJECT", rawJSON)
	if err != nil {
		return err
	}

	return nil
}

func cacheGeoFenceUsers(geoFenceId string, userId string) error {
	conn := database.RedisPool.Get()

	defer func() {
		err := conn.Close()
		if err != nil {
			panic(err)
		}
	}()

	key := NewGeoFenceUserId(geoFenceId)

	if _, err := conn.Do("SADD", key, userId); err != nil {
		return err
	}
	if _, err := conn.Do("EXPIRE", key, 60*3); err != nil {
		return err
	}

	return nil
}

func SetHookPoint(userId, discountId string, locations []discounts.Location) error {
	for _, l := range locations {
		keyId := newGeoFenceId(discountId, l)
		if err := setPoint(HookPointKey, keyId, l); err != nil {
			return err
		}

		if err := cacheGeoFenceUsers(keyId, userId); err != nil {
			return err
		}
	}

	return nil
}

func SetLocationPoint(discountId string, locations []discounts.Location) error {
	for _, l := range locations {
		if err := setPoint(PointKey, newLocationId(discountId, l), l); err != nil {
			return err
		}
	}
	return nil
}

func SetLocationHooks(discountId string, locations []discounts.Location) error {
	conn := database.GeoRedisPool.Get()

	defer func() {
		err := conn.Close()
		if err != nil {
			panic(err)
		}
	}()

	for _, location := range locations {
		nearby, err := nearbyLocations(location)

		if len(nearby) > 1 { // we've already hook which covered this case
			continue
		}

		endpoint := infra.GetEnv("DOCKER_APP_URL") + "/geo-fencing/nearby-hook"

		_, err = conn.Do(
			"SETHOOK",
			newGeoFenceId(discountId, location),
			endpoint,
			"NEARBY", HookPointKey, "FENCE", "POINT", location.Lat, location.Lon, nearbyLocation,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func DeleteGeoFence(discountId string, locations []discounts.Location) error {
	conn := database.GeoRedisPool.Get()

	defer func() {
		err := conn.Close()
		if err != nil {
			panic(err)
		}
	}()

	for _, location := range locations {
		geoFenceId := newGeoFenceId(discountId, location)

		if _, err := conn.Do("DELHOOK", geoFenceId); err != nil {
			return err
		}
		if _, err := conn.Do("DEL", HookPointKey, geoFenceId); err != nil {
			return err
		}
	}
	return nil
}

func nearbyLocations(searchLocation discounts.Location) ([]interface{}, error) {
	conn := database.GeoRedisPool.Get()

	defer func() {
		err := conn.Close()
		if err != nil {
			panic(err)
		}
	}()

	values, err := redis.Values(
		conn.Do("NEARBY", PointKey, "POINT", searchLocation.Lat, searchLocation.Lon, nearbyLocation),
	)

	if err != nil {
		return nil, err
	}

	if len(values) < 2 {
		return nil, err
	}

	// 0 - cursor
	values, err = redis.Values(values[1], nil)
	if err != nil {
		return nil, err
	}

	return values, nil
}

func CloseLocations(discountId string, searchLocation discounts.Location, discountLocations []discounts.Location) []discounts.Location {
	var locations []discounts.Location
	values, err := nearbyLocations(searchLocation)

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
			if newLocationId(discountId, discountLocation) != locationId {
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
