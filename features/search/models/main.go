package search

import (
	"github.com/zdunecki/discountly/features/discounts/models"
)

type Search struct {
	Keywords []string           `bson:"keywords,omitempty" json:"keywords,omitempty"`
	Location discounts.Location `bson:"location,omitempty" json:"location,omitempty"`
}
