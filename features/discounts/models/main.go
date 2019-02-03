package discounts

type PromoCodeStatus int

const (
	Inactive PromoCodeStatus = iota + 1
	Active
	Waiting
)

type Date struct {
	StartDate string `bson:"start_date,omitempty" json:"start_date,omitempty"`
	EndDate   string `bson:"end_date,omitempty" json:"end_date,omitempty"`
	Date      string `bson:"date,omitempty" json:"date,omitempty"`
}

type Discount struct {
	Id         string      `bson:"id,omitempty" json:"id,omitempty"`
	Name       string      `bson:"name,omitempty" json:"name,omitempty"`
	Keywords   []string    `bson:"keywords,omitempty" json:"keywords,omitempty"`
	ImageUrl   string      `bson:"image_url" json:"image_url"`
	Locations  []Location  `bson:"locations,omitempty" json:"locations,omitempty"`
	Rules      []Rule      `bson:"rules,omitempty" json:"rules,omitempty"`
	PromoCodes []PromoCode `bson:"promo_codes,omitempty" json:"promo_codes,omitempty"`
}

type ProtectedDiscount struct {
	Id        string     `bson:"id,omitempty" json:"id,omitempty"`
	Name      string     `bson:"name,omitempty" json:"name,omitempty"`
	Locations []Location `bson:"locations,omitempty" json:"locations,omitempty"`
	ImageUrl  string     `bson:"image_url" json:"image_url"`
}

type DiscountDefinition struct {
	UserId    string     `bson:"user_id" json:"user_id"`
	Company   string     `bson:"company" json:"company"`
	Discounts []Discount `bson:"discounts" json:"discounts"`
}


