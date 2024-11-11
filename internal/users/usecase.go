package users

import "bookloop.net/internal/models"

type UseCase interface {
	Register(input models.InputRegisterUser) (*models.User, error)
	Activate(tokenPlaintext string) (*models.User, error)
	Login(input models.InputLoginUser) (*models.User, error)
	NewActivationToken(input models.InputNewActivationTokenUser) (*models.User, error)
	Authenticate(token string) (*models.User, error)
}
