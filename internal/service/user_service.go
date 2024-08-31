package service

import "github.com/alwindoss/casa/internal/repository"

type UserService interface {
}

func NewUserService(repo repository.UserRepositoy) UserService {
	return &userService{}
}

type userService struct {
}
