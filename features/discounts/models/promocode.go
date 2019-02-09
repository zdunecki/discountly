package discounts

import (
	"github.com/zdunecki/discountly/lib"
	"math/rand"
	"strings"
	"time"
)

type PromoCode struct {
	StartDate string          `bson:"start_date,omitempty" json:"start_date,omitempty"`
	EndDate   string          `bson:"end_date,omitempty" json:"end_date,omitempty"`
	Status    PromoCodeStatus `bson:"status,omitempty" json:"status,omitempty"`
	Code      string          `bson:"code,omitempty" json:"code"`
}

func generateRandomCode() string {
	rand.Seed(time.Now().UnixNano())
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	length := 8
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	return b.String()
}

func (p *PromoCode) New() PromoCode {
	now := lib.Moment(nil)
	return PromoCode{
		StartDate: now.ISO(),
		EndDate:   now.Add(5, "minutes").ISO(),
		Status:    Active,
		Code: generateRandomCode(),
	}
}
