package discounts

import "github.com/satori/go.uuid"

type Location struct {
	Id  string `bson:"id,omitempty" json:"id,omitempty"`
	Lat float64 `bson:"lat,omitempty" json:"lat,omitempty"`
	Lon float64 `bson:"lon,omitempty" json:"lon,omitempty"`
}

func (location Location) New() Location {
	id := uuid.NewV4()
	return Location{
		Id:  id.String(),
		Lat: location.Lat,
		Lon: location.Lon,
	}
}

func (location Location) Edit() Location {
	return Location{
		Lat: location.Lat,
		Lon: location.Lon,
	}
}
