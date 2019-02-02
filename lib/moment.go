package lib

import (
	"errors"
	"github.com/jinzhu/now"
	"time"
)

type IMoment struct {
	day time.Time
	err error
}

func Moment(i interface{}) IMoment {
	var t time.Time
	switch v := i.(type) {
	case string: //TODO: find better solution how to handle this
		n, err := now.Parse(v)
		if err == nil {
			t = n
			break
		}
		n, err2 := time.Parse(time.RFC3339, v)
		if err2 != nil {
			panic(err)
		}
		t = n
	case time.Time:
		t = v
	case IMoment:
		t = v.day
	case error:
		panic(v)
	default:
		t = time.Now().Local()
	}

	return IMoment{
		t,
		nil,
	}
}
func (m IMoment) Add(n int, t string) IMoment {
	switch t {
	case "years":
		return Moment(m.day.AddDate(n, 0, 0))
	case "months":
		return Moment(m.day.AddDate(0, n, 0))
	case "days":
		return Moment(m.day.AddDate(0, 0, n))
	case "hours":
		return Moment(m.day.Add(time.Hour * (time.Duration(n))))
	case "minutes":
		return Moment(m.day.Add(time.Minute * (time.Duration(n))))
	default:
		err := errors.New("type not found")
		return Moment(err)
	}
}

func (m IMoment) ISO() string {
	return m.day.Format(time.RFC3339)
}

func (m IMoment) IsValid() bool {
	return m.err == nil
}

func (m IMoment) IsBetween(date1 string, date2 string) bool {
	c1, err := time.Parse(time.RFC3339, date1)
	if err != nil {
		panic(err)
	}
	c2, err := time.Parse(time.RFC3339, date2)
	if err != nil {
		panic(err)
	}

	return (m.day.After(c1) && m.day.Before(c2)) || m.day.Equal(c1) || m.day.Equal(c2)
}

func (m IMoment) Ymd() string {
	return m.day.Format("2006-01-02")
}

func (m IMoment) Weekday() int {
	return int(m.day.Weekday())
}

func (m IMoment) Native() time.Time {
	return m.day
}
