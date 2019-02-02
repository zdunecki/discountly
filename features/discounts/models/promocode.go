package discounts

import "github.com/zdunecki/discountly/lib"

type PromoCode struct {
	StartDate string          `bson:"start_date,omitempty" json:"start_date,omitempty"`
	EndDate   string          `bson:"end_date,omitempty" json:"end_date,omitempty"`
	Status    PromoCodeStatus `bson:"status,omitempty" json:"status,omitempty"`
}

func (p *PromoCode) New() PromoCode {
	now := lib.Moment(nil)
	return PromoCode{
		StartDate: now.ISO(),
		EndDate: now.Add(5, "minutes").ISO(),
		Status: Active,
	}
}