package notifications

import (
	"encoding/json"
	"github.com/gomodule/redigo/redis"
	"github.com/pusher/pusher-http-go"
	"github.com/zdunecki/discountly/db"
	_discounts "github.com/zdunecki/discountly/features/discounts"
	"github.com/zdunecki/discountly/features/discounts/models"
	"github.com/zdunecki/discountly/features/notifications/models"
	"github.com/zdunecki/discountly/features/search/finder"
	"log"
	"strings"
)

func findClosestDiscounts(conn redis.Conn, repo *database.Repository, userGeoFencing models.UserGeoFencing) []discounts.Discount {
	var closestFromCache []discounts.Discount
	cacheKey := "NOTIFICATION_RECIVE_HOOK_nearbyDiscounts"
	cached, _ := conn.Do("GET", cacheKey)

	if cached != nil {
		if err := json.Unmarshal(cached.([]byte), &closestFromCache); err != nil {
			log.Fatal(err)
		}
		return closestFromCache
	} else {
		values, err := finder.NearbyLocations(userGeoFencing.Location)

		if err != nil {
			log.Fatal(err)
		}

		var nearbyDiscountsIds []string

		for _, val := range values {
			v, err := redis.Strings(val, nil)
			if err != nil {
				log.Fatal(err)
			}

			keyId := v[0]
			discountId := strings.Split(keyId, "|")[0]

			nearbyDiscountsIds = append(nearbyDiscountsIds, discountId)
		}

		nearbyDiscounts, err := repo.Discounts.FindInIds(nearbyDiscountsIds)
		if err != nil {
			log.Fatal(err)
		}

		jDiscounts, err := json.Marshal(nearbyDiscounts)
		// TODO: EX is not the best solution, how to handle it better?
		if _, err := conn.Do("SET", cacheKey, jDiscounts, "EX", "60"); err != nil {
			log.Fatal(err)
		}
		return nearbyDiscounts
	}
}

func receiveHook(webHook *pusher.Webhook) {
	repo, err := database.NewRepo(database.DbConnection())

	if err != nil {
		log.Fatal(err)
	}
	defer repo.Discounts.Close()

	conn := database.RedisPool.Get()
	cache := make(map[string]bool)

	var cachedClosestDiscounts []discounts.Discount

	for _, event := range webHook.Events {
		if event.Channel != "private-connection" || event.Event != "client-geofencing" {
			continue
		}

		var userGeoFencing models.UserGeoFencing

		if err := json.Unmarshal([]byte(event.Data), &userGeoFencing); err != nil {
			log.Fatal(err)
		}

		if cachedClosestDiscounts == nil {
			cachedClosestDiscounts = findClosestDiscounts(conn, repo, userGeoFencing)
		}

		for _, discount := range cachedClosestDiscounts {
			if _, ok := cache[discount.Id]; !ok {

				if finder.ApplyTheRules(discount.Rules) {
					cache[discount.Id] = true
				} else {
					err := finder.DeleteGeoFence(
						discount.Id,
						[]discounts.Location{userGeoFencing.Location},
					)
					if err != nil {
						log.Fatal(err)
					}

					errorMessage := struct {
						Code       int    `json:"code"`
						DiscountId string `json:"discount_id"`
						Message    string `json:"message"`
					}{
						_discounts.NotApplyingToRulesCode,
						discount.Id,
						_discounts.NotApplyingToRulesMessage,
					}

					data := map[string]interface{}{"error": errorMessage}

					if _, err := Client.Trigger("private-"+userGeoFencing.UniqueRandomId, finder.HookPointKey, data); err != nil {
						log.Fatal(err)
					}

					log.Println("discount not applying to the rules")
					continue
				}
			}

			if err := finder.SetHookPoint(
				userGeoFencing.UniqueRandomId,
				discount.Id,
				[]discounts.Location{userGeoFencing.Location},
			); err != nil {
				log.Fatal(err)
			}
		}
	}
}
