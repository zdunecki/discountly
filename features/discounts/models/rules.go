package discounts

import (
	"github.com/zdunecki/discountly/lib"
	"github.com/satori/go.uuid"
)

type Rule struct {
	Id        string     `bson:"id,omitempty" json:"id,omitempty"`
	StartDate string     `bson:"start_date,omitempty" json:"start_date,omitempty"`
	EndDate   string     `bson:"end_date,omitempty" json:"end_date,omitempty"`
	Dates     []Date     `bson:"dates,omitempty" json:"dates,omitempty"`
	WeekDays  []WeekDays `bson:"week_days,omitempty" json:"week_days,omitempty"`
}

func (rule Rule) New() (Rule, error) {
	id := uuid.NewV4()
	createRule := Rule{
		Id:        id.String(),
		StartDate: lib.Moment(rule.StartDate).ISO(),
		EndDate:   lib.Moment(rule.EndDate).ISO(),
		Dates:     rule.dates(),
		WeekDays:  rule.weekdays(),
	}
	err := createRule.isInvalid()

	if err != nil {
		return Rule{}, err
	}

	return createRule, nil
}

func (rule Rule) Edit() (Rule, error) {
	createRule := Rule{
		StartDate: lib.Moment(rule.StartDate).ISO(),
		EndDate:   lib.Moment(rule.EndDate).ISO(),
		Dates:     rule.dates(),
		WeekDays:  rule.weekdays(),
	}
	err := createRule.isInvalid()

	if err != nil {
		return Rule{}, err
	}

	return createRule, nil
}

func (rule Rule) isInvalid() error {
	for _, weekDay := range rule.WeekDays {
		if err := weekDay.validWeekDays(); err != nil {
			return err
		}
	}

	return nil
}

func (rule Rule) dates() []Date {
	var newDates []Date

	for _, date := range rule.Dates {
		newDates = append(newDates, Date{
			StartDate: lib.Moment(date.StartDate).ISO(),
			EndDate:   lib.Moment(date.EndDate).ISO(),
			Date:      lib.Moment(date.StartDate).Ymd(),
		})
	}
	return newDates
}

func (rule Rule) weekdays() []WeekDays {
	var newDates []WeekDays

	for _, weekDay := range rule.WeekDays {
		var weekDays WeekDays
		if weekDay.StartDate != "" && weekDay.EndDate != "" {
			weekDays = WeekDays{
				StartDate: lib.Moment(weekDay.StartDate).ISO(),
				EndDate:   lib.Moment(weekDay.EndDate).ISO(),
				Day:       lib.Moment(weekDay.StartDate).Weekday(),
			}
		} else {
			weekDays = WeekDays{
				StartDate: "",
				EndDate:   "",
				Day:       weekDay.Day,
			}
		}
		newDates = append(newDates, weekDays)
	}
	return newDates
}
