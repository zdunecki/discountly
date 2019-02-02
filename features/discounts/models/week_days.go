package discounts

import (
	"errors"
	"github.com/zdunecki/discountly/lib"
)

type WeekDays struct {
	StartDate string `bson:"start_date,omitempty" json:"start_date,omitempty"`
	EndDate   string `bson:"end_date,omitempty" json:"end_date,omitempty"`
	Day       int    `bson:"day,omitempty" json:"day,omitempty"`
}

func (weekday WeekDays) validWeekDays() error {
	if lib.Moment(weekday.StartDate).Ymd() != lib.Moment(weekday.EndDate).Ymd() {
		return errors.New("start date and end date is different")
	}

	if lib.Moment(weekday.StartDate).Weekday() != weekday.Day {
		return errors.New("date week day is not the same as day")
	}

	return nil
}