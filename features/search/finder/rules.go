package finder

import (
	"github.com/gomodule/redigo/redis"
	"github.com/jinzhu/now"
	"github.com/zdunecki/discountly/features/discounts/models"
	"github.com/zdunecki/discountly/lib"
	"strconv"
	"time"
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

func ApplyTheRules(rules []discounts.Rule) bool {
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

/*
	TODO: check if caching to redis is better
*/
func ApplyTheRulesWithCaching(conn redis.Conn, discountId string, rules []discounts.Rule) (bool, error) {
	response, err := conn.Do("GET", "cache-discount-rules:"+discountId)

	if err != nil {
		return false, err
	}

	if response != nil {
		return true, nil
	}

	for _, rule := range rules {
		betweenStartEnd := lib.Moment(nil).IsBetween(rule.StartDate, rule.EndDate)
		inWeekDays := isInWeekDays(rule.WeekDays)
		absoluteDates := isInAbsoluteDates(rule.Dates)

		if !((betweenStartEnd && inWeekDays) || absoluteDates) {
			return false, nil
		}
	}

	ttl := strconv.Itoa(minimumRuleOfTheDay(rules))

	if _, err := conn.Do("SET", "cache-discount-rules:"+discountId, discountId, "EX", ttl); err != nil {
		return false, err
	}

	return true, nil
}

func minimumRuleOfTheDay(rules []discounts.Rule) int {
	var max string

	for _, rule := range rules {
		if max == "" || rule.EndDate > max {
			max = rule.EndDate
		}

		for _, w := range rule.WeekDays {
			if w.EndDate > max {
				max = w.EndDate
			}
		}
		for _, d := range rule.Dates {
			if d.EndDate > max {
				max = d.EndDate
			}
		}
	}

	var diff time.Duration

	maximumIsAfterThanToday := lib.Moment(max).Native().After(now.EndOfDay())

	if maximumIsAfterThanToday {
		lastDayHours := now.EndOfDay().Sub(time.Now())
		diff = lastDayHours
	} else {
		diff = lib.Moment(max).Native().Sub(time.Now())
	}

	return int(diff.Seconds())
}