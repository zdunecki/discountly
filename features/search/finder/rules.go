package finder

import (
	"github.com/zdunecki/discountly/features/discounts/models"
	"github.com/zdunecki/discountly/lib"
)

func isInWeekDays(weekDays []discounts.WeekDays) bool {
	if len(weekDays) == 0 {
		return true
	}

	now := lib.Moment(nil)

	for _, weekDay := range weekDays {
		if now.Weekday() == weekDay.Day {
			if weekDay.StartDate != "" && weekDay.EndDate != "" {
				return now.IsBetween(weekDay.StartDate, weekDay.EndDate)
			}
			return true
		}
	}

	return false
}

func isInAbsoluteDates(dates []discounts.Date) bool {
	if dates == nil {
		return false
	}
	if len(dates) == 0 {
		return true
	}

	now := lib.Moment(nil)

	for _, weekDay := range dates {
		if now.Ymd() == weekDay.Date {
			return now.IsBetween(weekDay.StartDate, weekDay.EndDate)
		}
	}

	return false
}

func applyTheRules(rules []discounts.Rule) bool {
	someInRules := func() bool {
		for _, rule := range rules {
			betweenStartEnd := lib.Moment(nil).IsBetween(rule.StartDate, rule.EndDate)
			inWeekDays := isInWeekDays(rule.WeekDays)
			absoluteDates := isInAbsoluteDates(rule.Dates)

			if (betweenStartEnd && inWeekDays) || absoluteDates {
				return true
			}
		}
		return false
	}
	return someInRules()
}
