package usecase

import (
	"errors"
	"nazartaraniuk/alertsProject/internal/domain"
	"nazartaraniuk/alertsProject/internal/repository"
)

type UserService interface {
	LoginUser(req domain.LoginReq) (*domain.User, error)
	RegisterUser(req domain.RegisterReq) (int64, error)
}

type UserServiceImpl struct {
	repository *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) UserService {
	return &UserServiceImpl{
		repository: repo,
	}
}

func (s *UserServiceImpl) LoginUser(req domain.LoginReq) (*domain.User, error) {
	user, err := s.repository.GetUserByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (s *UserServiceImpl) RegisterUser(req domain.RegisterReq) (int64, error) {
	id, err := s.repository.CreateUser(req)
	if err != nil {
		return id, err
	}
	return id, nil
}
