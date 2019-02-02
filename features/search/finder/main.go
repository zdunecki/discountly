package finder

import (
	"github.com/zdunecki/discountly/features/discounts/models"
	"github.com/zdunecki/discountly/features/search/models"
)

func FindBestDiscounts(definitions []discounts.DiscountDefinition, search search.Search) []discounts.Discount {
	var result []discounts.Discount

	for _, definition := range definitions {
		for _, discount := range definition.Discounts {
			if !hasKeyword(discount.Keywords, search) || !applyTheRules(discount.Rules) {
				continue
			}

			locations := CloseLocations(search, discount.Locations)
			if locations == nil {
				continue
			}

			result = append(
				result,
				discountByBestLocations(discount, locations),
			)
		}
	}

	if len(result) == 0 {
		return nil
	} else {
		return result
	}
}
