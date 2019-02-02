package auth

import "github.com/satori/go.uuid"

type User struct {
	Id    string `bson:"id" json:"id"`
	Email string `bson:"email" json:"email"`
	Name  string `bson:"name" json:"name"`
}

func (user User) New() User {
	id := uuid.NewV4()
	return User{
		Id:    id.String(),
		Email: user.Email,
		Name:  user.Name,
	}
}
