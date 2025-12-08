package user_test

import (
	mock_notification "belajarGo2/service/notification/mock"
	"belajarGo2/service/user"
	mock_user "belajarGo2/service/user/mock"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/pobyzaarif/goshortcute"
	"github.com/stretchr/testify/assert"
)

var loggerOption = slog.HandlerOptions{AddSource: true}
var logger = slog.New(slog.NewJSONHandler(os.Stdout, &loggerOption))

func TestRegister(t *testing.T) {
	tests := []struct {
		name      string
		inputUser user.User
		mockUser  func(m *mock_user.MockRepository)
		mockNotif func(m *mock_notification.MockRepository)
		wantErr   bool
	}{
		{
			name:      "error on GetByEmail",
			inputUser: user.User{},
			mockUser: func(m *mock_user.MockRepository) {
				m.EXPECT().GetByEmail("").Return(user.User{}, errors.New("record not found"))
			},
			mockNotif: func(m *mock_notification.MockRepository) {},
			wantErr:   true,
		},
		{
			name:      "error email already exists",
			inputUser: user.User{Email: "test@example.com"},
			mockUser: func(m *mock_user.MockRepository) {
				m.EXPECT().GetByEmail("test@example.com").Return(user.User{ID: "1", Email: "test@example.com"}, nil)
			},
			mockNotif: func(m *mock_notification.MockRepository) {},
			wantErr:   true,
		},
		{
			name:      "error when created user",
			inputUser: user.User{Email: "test@example.com"},
			mockUser: func(m *mock_user.MockRepository) {
				m.EXPECT().GetByEmail("test@example.com").Return(user.User{}, nil)
				m.EXPECT().Create(gomock.Any()).Return(errors.New("db con error"))
			},
			mockNotif: func(m *mock_notification.MockRepository) {},
			wantErr:   true,
		},
		{
			name:      "success",
			inputUser: user.User{Email: "test@example.com"},
			mockUser: func(m *mock_user.MockRepository) {
				m.EXPECT().GetByEmail("test@example.com").Return(user.User{}, nil)
				m.EXPECT().Create(gomock.Any()).Return(nil)
			},
			mockNotif: func(m *mock_notification.MockRepository) {
				m.EXPECT().SendEmail("", "test@example.com", gomock.Any(), gomock.Any()).Return(nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mock_userRepo := mock_user.NewMockRepository(ctrl)
			mock_notification := mock_notification.NewMockRepository(ctrl)

			tt.mockUser(mock_userRepo)
			tt.mockNotif(mock_notification)

			productService := user.NewService(
				logger,
				mock_userRepo,
				"http://appDeploymentUrl.com",
				"exampleexampleexampleexampleexampleexampleexampleexampleexampleexample",
				"32character32character32characte",
				mock_notification,
			)

			id, err := productService.Register(tt.inputUser)
			if tt.wantErr {
				assert.Equal(t, "", id)
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestVerifyEmail(t *testing.T) {
	key := []byte("32character32character32characte")
	tsInTheFuture := time.Now().Add(time.Minute * 10).Unix()
	notExpiredCode, _ := goshortcute.AESCBCEncrypt([]byte(fmt.Sprintf("%s|%d", "email@mail.com", tsInTheFuture)), key)

	tests := []struct {
		name      string
		inputUser string
		mockUser  func(m *mock_user.MockRepository)
		mockNotif func(m *mock_notification.MockRepository)
		wantErr   bool
	}{
		{
			name:      "error invalid decrypt",
			inputUser: "dhslkashdlaskdh",
			mockUser:  func(m *mock_user.MockRepository) {},
			mockNotif: func(m *mock_notification.MockRepository) {},
			wantErr:   true,
		},
		{
			name:      "error invalid text format",
			inputUser: "DNoioYCcyLIiSFnk2K86BA==",
			mockUser:  func(m *mock_user.MockRepository) {},
			mockNotif: func(m *mock_notification.MockRepository) {},
			wantErr:   true,
		},
		{
			name:      "error invalid ts format",
			inputUser: "DUvjHjyBKrkZ2k88+Pkj6g==",
			mockUser:  func(m *mock_user.MockRepository) {},
			mockNotif: func(m *mock_notification.MockRepository) {},
			wantErr:   true,
		},
		{
			name:      "error invalid ts format",
			inputUser: "DUvjHjyBKrkZ2k88+Pkj6g==",
			mockUser:  func(m *mock_user.MockRepository) {},
			mockNotif: func(m *mock_notification.MockRepository) {},
			wantErr:   true,
		},
		{
			name:      "error invalid ts expired",
			inputUser: "bkCPB5TdD40yVXwL+jGnjDpwRFE7nZ4vj0Aw3N3nkwQ=",
			mockUser:  func(m *mock_user.MockRepository) {},
			mockNotif: func(m *mock_notification.MockRepository) {},
			wantErr:   true,
		},
		{
			name:      "error ts not expired but get by email error",
			inputUser: notExpiredCode,
			mockUser: func(m *mock_user.MockRepository) {
				m.EXPECT().GetByEmail("email@mail.com").Return(user.User{}, errors.New("db error"))
			},
			mockNotif: func(m *mock_notification.MockRepository) {},
			wantErr:   true,
		},
		{
			name:      "error ts not expired but get by email succes but the email already verified",
			inputUser: notExpiredCode,
			mockUser: func(m *mock_user.MockRepository) {
				m.EXPECT().GetByEmail("email@mail.com").Return(user.User{IsEmailVerified: true}, nil)
			},
			mockNotif: func(m *mock_notification.MockRepository) {},
			wantErr:   true,
		},
		{
			name:      "error ts not expired but get by email succes but error when update email verification",
			inputUser: notExpiredCode,
			mockUser: func(m *mock_user.MockRepository) {
				m.EXPECT().GetByEmail("email@mail.com").Return(user.User{}, nil)
				m.EXPECT().UpdateEmailVerification(user.User{IsEmailVerified: true}).Return(errors.New("db error"))
			},
			mockNotif: func(m *mock_notification.MockRepository) {},
			wantErr:   true,
		},
		{
			name:      "success",
			inputUser: notExpiredCode,
			mockUser: func(m *mock_user.MockRepository) {
				m.EXPECT().GetByEmail("email@mail.com").Return(user.User{}, nil)
				m.EXPECT().UpdateEmailVerification(user.User{IsEmailVerified: true}).Return(nil)
			},
			mockNotif: func(m *mock_notification.MockRepository) {},
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mock_userRepo := mock_user.NewMockRepository(ctrl)
			mock_notification := mock_notification.NewMockRepository(ctrl)

			tt.mockUser(mock_userRepo)
			tt.mockNotif(mock_notification)

			productService := user.NewService(
				logger,
				mock_userRepo,
				"http://appDeploymentUrl.com",
				"exampleexampleexampleexampleexampleexampleexampleexampleexampleexample",
				string(key),
				mock_notification,
			)

			err := productService.VerifyEmail(tt.inputUser)
			if tt.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
