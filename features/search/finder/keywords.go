package finder

import (
	"github.com/zdunecki/discountly/features/search/models"
	"github.com/zdunecki/discountly/lib"
)

func hasKeyword(keywords []string, search search.Search) bool {
	for _, keyword := range search.Keywords {
		if lib.Contains(keywords, keyword) {
			return true
		}
	}
	return false
}
