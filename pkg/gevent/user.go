package gevent

import "github.com/common-fate/common-fate/pkg/identity"

type User struct {
	ID        string `json:"id" dynamodbav:"id"`
	Email     string `json:"email" dynamodbav:"email"`
	FirstName string `json:"firstName" dynamodbav:"firstName"`
	LastName  string `json:"lastName" dynamodbav:"lastName"`
}

func UserFromIdentityUser(user identity.User) User {
	return User{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}
}
