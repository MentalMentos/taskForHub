package converter

import (
	"github.com/MentalMentos/taskForHub/auth/internal/model"
)

func ToApi(u model.User) model.UserApi {
	return model.UserApi{
		ID:    u.ID.Hex(),
		Name:  u.Name,
		Email: u.Email,
	}
}
