package e2e

import (
	"github.com/bouk/monkey"
	"github.com/stretchr/testify/assert"
	"github.com/zdunecki/discountly/db"
	"github.com/zdunecki/discountly/features/auth/models"
	"github.com/zdunecki/discountly/features/discounts/models"
	"github.com/zdunecki/discountly/features/search/finder"
	"github.com/zdunecki/discountly/features/search/models"
	"github.com/zdunecki/discountly/lib"
	"github.com/zdunecki/discountly/tests"
	"testing"
	"time"
)

var testLocations = []discounts.Location{
	{
		Lat: testLat,
		Lon: testLon,
	},
}

func createDiscounts() []discounts.Discount {
	december2018th27 := lib.Moment(nil)
	fiveHoursAfter := december2018th27.Add(2, "hours")

	twoDaysBefore := december2018th27.Add(-2, "days")
	twoDaysAfter := december2018th27.Add(2, "days")

	twoDaysAfterAnd5hAfter := twoDaysAfter.Add(5, "hours")
	twoDaysAfterAnd10hAfter := twoDaysAfter.Add(10, "hours")

	return []discounts.Discount{
		{
			Name:       "test-name",
			Keywords:   []string{"abc"},
			Locations:  testLocations,
			PromoCodes: []discounts.PromoCode{},
			ImageUrl:   "image-url",
			Rules: []discounts.Rule{
				{
					StartDate: twoDaysBefore.ISO(),
					EndDate:   twoDaysAfter.ISO(),
					Dates: []discounts.Date{
						{
							StartDate: twoDaysAfterAnd5hAfter.ISO(),
							EndDate:   twoDaysAfterAnd10hAfter.ISO(),
						},
					},
					WeekDays: []discounts.WeekDays{
						{
							StartDate: december2018th27.ISO(),
							EndDate:   fiveHoursAfter.ISO(),
						},
					},
				},
			},
		},
	}
}

func TestFindBestDiscounts(t *testing.T) {
	assert := assert.New(t)

	wayback := time.Date(2018, time.December, 27, 0, 0, 0, 0, time.UTC)
	patch := monkey.Patch(time.Now, func() time.Time { return wayback })
	defer patch.Unpatch()

	repo, err := database.NewRepo(dbAddress)
	if err != nil {
		panic(err)
	}

	defer repo.Discounts.Close()
	defer tests.DeleteAll()

	userId := "test-user-id"

	_, err = repo.Discounts.CreateDefinition(auth.User{
		Id:    userId,
		Email: "test-email@gmail.com",
		Name:  "test",
	})

	createdDiscount, _ := repo.Discounts.CreateDiscounts(
		userId,
		createDiscounts(),
	)

	_ = finder.SetLocationPoint(createdDiscount[0].Id, createdDiscount[0].Locations)

	criteria := search.Search{
		Keywords: []string{"abc"},
		Location: discounts.Location{
			Lat: nearbyTestLat,
			Lon: nearbyTestLon,
		},
	}

	allByKeyWords, _ := repo.Discounts.FindAllByKeywords(criteria)

	result := finder.FindBestDiscounts(allByKeyWords, criteria)

	assert.Equal(len(result), 1)
}
