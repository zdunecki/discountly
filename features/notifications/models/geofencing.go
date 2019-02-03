package models

import "github.com/zdunecki/discountly/features/discounts/models"

type UserGeoFencing struct {
	UniqueRandomId string             `json:"unique_random_id"`
	Location       discounts.Location `json:"location"`
}
