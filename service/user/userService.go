package user

import (
	"belajarGo2/service/notification"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/pobyzaarif/goshortcute"
	"golang.org/x/crypto/bcrypt"
)

type service struct {
	logger                  *slog.Logger
	repo                    Repository
	appDeploymentUrl        string
	jwtSign                 string
	appEmailVerificationKey string
	notifRepo               notification.Repository
}

const (
	verificationCodeTTL = 5
)

type Service interface {
	Register(user User) (id string, err error)
	Login(username string, password string) (accessToken string, err error)
	GetByEmail(email string) (user User, err error)
	VerifyEmail(verificationCodeEncrypt string) (err error)
}

func NewService(logger *slog.Logger, repo Repository, appDeploymentUrl string, jwtSign string, appEmailVerificationKey string, notifRepo notification.Repository) Service {
	return &service{
		logger:                  logger,
		repo:                    repo,
		appDeploymentUrl:        appDeploymentUrl,
		jwtSign:                 jwtSign,
		appEmailVerificationKey: appEmailVerificationKey,
		notifRepo:               notifRepo,
	}
}

const (
	SubjectRegisterAccount = "Activate Your Account!"
	// EmailBodyRegisterAccount = `Halo, %v, Aktivasi akun anda dengan membuka tautan dibawah<br><br/>%v`
	EmailBodyRegisterAccount = `Halo, %v, Aktivasi akun anda dengan membuka tautan dibawah<br><br/>%v<br/>catatan: link hanya berlaku %v menit`
)

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

	if err = s.repo.Create(user); err != nil {
		return
	}

	// activationLink := "http://localhost:8080/users/verify-email"
	timeNow := time.Now()
	expAt := timeNow.Add(time.Duration(time.Minute * verificationCodeTTL)).Unix()

	verificationCode := fmt.Sprintf("%v|%v", user.Email, expAt)
	verificationCodeEncrypt, _ := goshortcute.AESCBCEncrypt([]byte(verificationCode), []byte(s.appEmailVerificationKey))
	// activationLink := s.appDeploymentUrl + "/users/email-verification/" + verificationCodeEncrypt
	verifCode := goshortcute.StringtoBase64Encode(verificationCodeEncrypt)
	activationLink := s.appDeploymentUrl + "/users/email-verification/" + verifCode

	_ = s.notifRepo.SendEmail(user.Fullname, user.Email, SubjectRegisterAccount, fmt.Sprintf(EmailBodyRegisterAccount, user.Fullname, activationLink, verificationCodeTTL))

	// Create user
	return user.ID, nil
}

func (s *service) VerifyEmail(verificationCodeEncrypt string) (err error) {
	// verificationCodeDecrypt, err := goshortcute.AESCBCDecrypt([]byte(verificationCodeEncrypt), []byte(s.appEmailVerificationKey))
	verifCodeDecode := goshortcute.StringtoBase64Decode(verificationCodeEncrypt)
	verificationCodeDecrypt, err := goshortcute.AESCBCDecrypt([]byte(verifCodeDecode), []byte(s.appEmailVerificationKey))
	if err != nil {
		s.logger.Error("verify email err", slog.Any("err", err.Error()))
		return errors.New("invalid or expired url")
	}

	verificationCode := strings.Split(verificationCodeDecrypt, "|")
	if len(verificationCode) != 2 {
		s.logger.Error("verify email err", slog.Any("err", verificationCodeDecrypt))
		return errors.New("invalid or expired url")
	}

	email := verificationCode[0]
	expAtStr := verificationCode[1]

	ts, err := strconv.ParseInt(expAtStr, 10, 64)
	if err != nil {
		s.logger.Error("verify email err", slog.Any("err", verificationCodeDecrypt))
		return errors.New("invalid or expired url")
	}
	expAt := time.Unix(ts, 0)
	if time.Now().After(expAt) {
		return errors.New("invalid or expired url")
	}

	getUser, err := s.repo.GetByEmail(email)
	if err != nil {
		s.logger.Error("verify email err", slog.Any("err", err))
		return err
	}

	if getUser.IsEmailVerified {
		s.logger.Error("verify email err", slog.Any("err", "email already verified"))
		return errors.New("invalid or expired url")
	}

	getUser.IsEmailVerified = true
	if err := s.repo.UpdateEmailVerification(getUser); err != nil {
		s.logger.Error("verify email err", slog.Any("err", err))
		return err
	}

	return nil
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

	if !getUser.IsEmailVerified {
		err = errors.New("email address has not been verified")
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
