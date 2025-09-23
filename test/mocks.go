package test

import "nazartaraniuk/alertsProject/internal/domain"

type mockUserService struct {
	RegisterUserFunc func(domain.RegisterReq) (int64, error)
	LoginUserFunc    func(domain.LoginReq) (*domain.User, error)
}

func (m *mockUserService) LoginUser(req domain.LoginReq) (*domain.User, error) {
	return m.LoginUserFunc(req)
}

func (m *mockUserService) RegisterUser(req domain.RegisterReq) (int64, error) {
	return m.RegisterUserFunc(req)
}


