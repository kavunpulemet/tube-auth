package mappers

import (
	dbmodels "auth/internal/database/models"
	"auth/internal/models"
)

func MapToDBUser(serviceUser models.User) dbmodels.User {
	return dbmodels.User{
		Id:       "",
		Email:    serviceUser.Email,
		Username: serviceUser.Username,
		Password: serviceUser.Password,
	}
}
