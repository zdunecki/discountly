package notifications

import (
	"encoding/json"
	"github.com/pusher/pusher-http-go"
	"github.com/zdunecki/discountly/db"
	_discounts "github.com/zdunecki/discountly/features/discounts"
	"github.com/zdunecki/discountly/features/discounts/models"
	"github.com/zdunecki/discountly/features/notifications/models"
	"github.com/zdunecki/discountly/features/search/finder"
	"log"
)

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

		//TODO: cache mechanism, pagination
		if cachedClosestDiscounts == nil {
			var closestFromRedis []discounts.Discount
			cached, _ := conn.Do("GET", "NOTIFICATION_RECIVE_HOOK_allDiscounts")

			if cached != nil {
				if err := json.Unmarshal(cached.([]byte), &closestFromRedis); err != nil {
					log.Fatal(err)
				}
				cachedClosestDiscounts = closestFromRedis
			} else {
				allDiscounts, err := repo.Discounts.FindAll()
				if err != nil {
					log.Fatal(err)
				}

				jDiscounts, err := json.Marshal(allDiscounts)
				// TODO: EX is not the best solution, how to handle it better?
				if _, err := conn.Do("SET", "NOTIFICATION_RECIVE_HOOK_allDiscounts", jDiscounts, "EX", "60"); err != nil {
					log.Fatal(err)
				}

				cachedClosestDiscounts = allDiscounts
			}

		}

		var userGeoFencing models.UserGeoFencing

		if err := json.Unmarshal([]byte(event.Data), &userGeoFencing); err != nil {
			log.Fatal(err)
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
