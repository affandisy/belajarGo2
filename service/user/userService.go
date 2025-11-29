package user

import (
	"errors"
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type service struct {
	logger  *slog.Logger
	repo    Repository
	jwtSign string
}

type Service interface {
	Register(user User) (id string, err error)
	Login(username string, password string) (accessToken string, err error)
	GetByEmail(email string) (user User, err error)
}

func NewService(logger *slog.Logger, repo Repository, jwtSign string) Service {
	return &service{
		logger:  logger,
		repo:    repo,
		jwtSign: jwtSign,
	}
}

func (s *service) Register(user User) (id string, err error) {
	// Find user by email
	getUser, err := s.repo.GetByEmail(user.Email)
	if err != nil {
		return
	}

	if getUser.Email != "" {
		err = errors.New("email registered already")
		return
	}

	// Hashing plain pass
	encPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return
	}

	user.ID = uuid.NewString()
	user.Password = string(encPassword)
	user.Role = "user"

	// Create user
	return "0", s.repo.Create(user)
}

func (s *service) Login(email string, password string) (accessToken string, err error) {
	getUser, err := s.repo.GetByEmail(email)
	if err != nil {
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(getUser.Password), []byte(password)); err != nil {
		s.logger.Error("login err", slog.Any("err", err.Error()))

		err = errors.New("wrong email or password")
		return "", err
	}

	token, err := s.generateToken(s.jwtSign, getUser.ID, getUser.Role)
	if err != nil {
		s.logger.Error("generate token err", slog.Any("err", err.Error()))

		err = errors.New("generate token error")
		return "", err
	}

	return token, err
}

func (s *service) GetByEmail(email string) (user User, err error) {
	return s.repo.GetByEmail(email)
}

func (s *service) generateToken(jwtSign string, id string, role string) (signedToken string, err error) {
	type jwtClaims struct {
		ID   string `json:"id"`
		Role string `json:"role"`
		jwt.RegisteredClaims
	}

	timeNow := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims{
		ID:   id,
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(timeNow),
			ExpiresAt: jwt.NewNumericDate(timeNow.Add(time.Hour * 24)),
		},
	})

	signedToken, err = token.SignedString([]byte(jwtSign))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
