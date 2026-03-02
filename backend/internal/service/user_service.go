package service

import (
	"log"
	"oct-backend/internal/model"
	"oct-backend/internal/repository"
	"oct-backend/internal/utils"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type UserService struct {
	Repo *repository.UserRepository
}

func (s *UserService) GetUserByID(id string) (*model.User, error) {
	return s.Repo.FindByID(id)
}

func (s *UserService) GetUserByUsername(username string) (*model.User, error) {
	return s.Repo.FindByUsername(username)
}

func (s *UserService) Register(username, password, email, gender string, age int) error {
	hash, err := utils.HashPassword(password)
	if err != nil {
		return err
	}
	user := &model.User{
		Username:  username,
		Password:  hash,
		Email:     email,
		Gender:    gender,
		Age:       age,
		CreatedAt: time.Now().Unix(),
	}
	return s.Repo.Create(user)
}

func (s *UserService) Login(username, password string) (string, error) {
	user, err := s.Repo.FindByUsername(username)
	if err != nil || !utils.CheckPasswordHash(password, user.Password) {
		return "", err
	}

	now := time.Now().Unix()
	if err := s.Repo.UpdateLastLoginAt(user.ID, now); err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  user.ID.Hex(),
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	})
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}
