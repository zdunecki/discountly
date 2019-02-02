package finder

import (
	"github.com/zdunecki/discountly/features/discounts/models"
)

//TODO: filter close locations with best for current user
func discountByBestLocations(d discounts.Discount, closestLocations []discounts.Location) discounts.Discount {
	return discounts.Discount{
		Id:         d.Id,
		Name:       d.Name,
		Keywords:   d.Keywords,
		Locations:  closestLocations,
		PromoCodes: d.PromoCodes,
		ImageUrl:   d.ImageUrl,
		Rules:      d.Rules,
	}
}
