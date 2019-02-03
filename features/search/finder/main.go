package finder

import (
	"github.com/zdunecki/discountly/features/discounts/models"
	"github.com/zdunecki/discountly/features/search/models"
)

func FindBestDiscounts(d []discounts.Discount, search search.Search) []discounts.Discount {
	var result []discounts.Discount

	for _, discount := range d {
		if !hasKeyword(discount.Keywords, search) || !ApplyTheRules(discount.Rules) {
			continue
		}

		locations := CloseLocations(discount.Id, search.Location, discount.Locations)

		if locations == nil {
			continue
		}

		result = append(
			result,
			discount,
		)
	}

	if len(result) == 0 {
		return nil
	} else {
		return result
	}
}