package usecase

import (
	"errors"
	"log/slog"

	"bookloop.net/config"
	"bookloop.net/internal/models"
	"bookloop.net/internal/users"
	"bookloop.net/internal/users/repository"
	"bookloop.net/pkg/db"
	"bookloop.net/pkg/validator"
)

type usersUC struct {
	cfg       *config.Config
	usersRepo users.Repository
	logger    *slog.Logger
}

func NewUsersUseCase(cfg *config.Config, usersRepo users.Repository, logger *slog.Logger) users.UseCase {
	return &usersUC{cfg: cfg, usersRepo: usersRepo, logger: logger}
}

func (u *usersUC) Register(input models.InputRegisterUser) (*models.User, error) {
	user := &models.User{
		Name:      input.Name,
		Email:     input.Email,
		Activated: false,
	}

	err := user.Password.Set(input.Password)
	if err != nil {
		return nil, err
	}

	v := validator.New()

	if models.ValidateUser(v, user); !v.Valid() {
		validationError := &validator.ValidationError{
			Errors: v.Errors,
			Err:    validator.ErrJSONIsNotValid,
		}
		return nil, validationError
	}

	err = u.usersRepo.Insert(user)
	if err != nil {
		if errors.Is(err, repository.ErrDuplicateEmail) {
			v.AddError("email", "a user with this email already exists")
			validationError := &validator.ValidationError{
				Errors: v.Errors,
				Err:    validator.ErrJSONIsNotValid,
			}
			return nil, validationError
		}
		return nil, err
	}

	return user, nil
}

func (u *usersUC) Activate(tokenPlaintext string) (*models.User, error) {
	v := validator.New()

	if models.ValidateTokenPlainText(v, tokenPlaintext); !v.Valid() {
		validationError := &validator.ValidationError{
			Errors: v.Errors,
			Err:    validator.ErrJSONIsNotValid,
		}
		return nil, validationError
	}

	user, err := u.usersRepo.GetForToken(models.ScopeActivation, tokenPlaintext)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			v.AddError("token", "invalid or expired activation token")
			validationError := &validator.ValidationError{
				Errors: v.Errors,
				Err:    validator.ErrJSONIsNotValid,
			}
			return nil, validationError
		}
		return nil, err
	}

	user.Activated = true

	err = u.usersRepo.Update(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *usersUC) Login(input models.InputLoginUser) (*models.User, error) {
	v := validator.New()

	models.ValidateEmail(v, input.Email)
	models.ValidatePasswordPlaintext(v, input.Password)

	if !v.Valid() {
		validationError := &validator.ValidationError{
			Errors: v.Errors,
			Err:    validator.ErrJSONIsNotValid,
		}
		return nil, validationError
	}

	user, err := u.usersRepo.GetByEmail(input.Email)
	if err != nil {
		return nil, err
	}

	match, err := user.Password.Matches(input.Password)
	if err != nil {
		return nil, err
	}

	if !match {
		return nil, db.ErrRecordNotFound
	}

	return user, nil
}

func (u *usersUC) NewActivationToken(input models.InputNewActivationTokenUser) (*models.User, error) {
	v := validator.New()

	if models.ValidateEmail(v, input.Email); !v.Valid() {
		validationError := &validator.ValidationError{
			Errors: v.Errors,
			Err:    validator.ErrJSONIsNotValid,
		}
		return nil, validationError
	}

	user, err := u.usersRepo.GetByEmail(input.Email)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			v.AddError("email", "no matching email address found")
			validationError := &validator.ValidationError{
				Errors: v.Errors,
				Err:    validator.ErrJSONIsNotValid,
			}
			return nil, validationError
		}
		return nil, err
	}

	if user.Activated {
		v.AddError("email", "user has already been activated")
		validationError := &validator.ValidationError{
			Errors: v.Errors,
			Err:    validator.ErrJSONIsNotValid,
		}
		return nil, validationError
	}

	return user, err
}

func (u *usersUC) Authenticate(token string) (*models.User, error) {
	return u.usersRepo.GetForToken(models.ScopeAuthentication, token)
}
