package search

import (
	"github.com/zdunecki/discountly/features/discounts/models"
	"strings"
)

type SearchDiscount struct {
	discounts.ProtectedDiscount `bson:",inline"`
	Token                       string `bson:"token,omitempty" json:"token"`
}

type Search struct {
	Keywords []string           `bson:"keywords,omitempty" json:"keywords,omitempty"`
	Location discounts.Location `bson:"location,omitempty" json:"location,omitempty"`
}

type Hook struct {
	Detect string `json:"detect"`
	Key    string `json:"key"`
	Id     string `json:"id"`
}

type HookGeoFenceId struct {
	DiscountId string
}

func (h *Hook) Split() HookGeoFenceId {
	split := strings.Split(h.Id, "|")

	return HookGeoFenceId{
		split[0],
	}
}
