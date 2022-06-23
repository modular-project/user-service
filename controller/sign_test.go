package controller

import (
	"errors"
	"testing"
	"time"
	"users-service/mocks"
	"users-service/model"
	"users-service/pkg"

	"github.com/gbrlsnchs/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestValidateRefreshToken(t *testing.T) {
	tests := []struct {
		give    string
		WantID  uint
		WantErr error
	}{
		{
			give:    "OK",
			WantID:  1,
			WantErr: nil,
		}, {
			give:    "Bad",
			WantID:  0,
			WantErr: ErrInvalidRefreshToken,
		},
	}
	for _, tt := range tests {
		to := mocks.NewTokener(t)
		re := mocks.NewRefreshStorager(t)
		opt := jwt.Options{
			Public: map[string]interface{}{
				"id":  1,
				"fgp": tt.give,
				"uid": 1,
			},
		}
		none := jwt.None()
		ts, err := jwt.Sign(none, &opt)
		if err != nil {
			t.Fatalf("fatal at sign jwt: %s", err)
		}
		jt, _ := jwt.FromString(ts)
		to.On("Validate", mock.Anything).Return(jt, nil)
		re.On("Find", mock.Anything).Return(model.Refresh{Hash: "OK"}, nil)
		ss := NewSignService(re, nil, nil, to)
		gotID, gotErr := ss.validateRefreshToken(&tt.give)
		assert.Equal(t, gotID, tt.WantID)
		if !errors.Is(gotErr, tt.WantErr) {
			t.Errorf("got err: %s, want err: %s", gotErr, tt.WantErr)
		}
	}
}

func TestSingIn(t *testing.T) {
	tests := []struct {
		give        model.LogIn
		WantErr     bool
		WantToken   string
		WantRefresh string
	}{
		{
			give:        model.LogIn{User: "OK@mail.com", Password: "Password12345."},
			WantToken:   "TOKENUSER1",
			WantRefresh: "REFRESHUSER1",
		}, {
			give:    model.LogIn{User: "Bad"},
			WantErr: true,
		}, {
			give:    model.LogIn{User: "OK@mail.com", Password: "WrongPassword12345."},
			WantErr: true,
		},
	}
	assert := assert.New(t)
	si := mocks.NewSignStorager(t)
	to := mocks.NewTokener(t)
	re := mocks.NewRefreshStorager(t)
	va := NewUserValidate()
	pwd, err := hashAndSalt([]byte("Password12345."))
	assert.Nil(err)

	password := string(pwd)
	si.On("Find", "OK@mail.com").Return(model.LogIn{ID: 1, Password: password}, nil)
	//si.On("Create", mock.Anything).Return(nil)
	to.On("Create", uint(1), uint(pkg.USER)).Return("TOKENUSER1", nil)
	to.On("CreateRefresh", mock.Anything, mock.Anything, mock.Anything).Return("REFRESHUSER1", nil)
	re.On("Create", mock.Anything).Return(nil)
	sign := NewSignService(re, va, si, to)
	for _, tt := range tests {
		gotToken, gotRefresh, gotErr := sign.SignIn(&tt.give)
		assert.Equal(gotToken, tt.WantToken)
		assert.Equal(gotRefresh, tt.WantRefresh)
		if tt.WantErr {
			assert.NotNil(gotErr)
		} else {
			assert.Nil(gotErr)
		}
	}

}

func TestVerify(t *testing.T) {
	tests := []struct {
		giveUID  uint
		giveCode string
		wantErr  error
	}{
		{
			giveUID:  1,
			giveCode: "CODEBAD",
			wantErr:  ErrInvalidCode,
		}, {
			giveUID:  1,
			giveCode: "CODEOK",
		}, {
			giveUID:  2,
			giveCode: "CODEOK",
			wantErr:  pkg.ErrNoRowsAffected,
		}, {
			giveUID:  3,
			giveCode: "EXPIREDCODE",
			wantErr:  ErrExpiredCode,
		},
	}
	ust := mocks.NewUserStorager(t)
	ver := mocks.NewVerificationStorager(t)
	ver.On("Find", uint(1)).Return(model.Verification{Code: "CODEOK", ExpiresAt: time.Now().Add(1 * time.Minute)}, nil)
	ver.On("Find", uint(2)).Return(model.Verification{}, pkg.ErrNoRowsAffected)
	ver.On("Find", uint(3)).Return(model.Verification{Code: "EXPIREDCODE", ExpiresAt: time.Now().Add(-1 * time.Minute)}, nil)
	ver.On("Delete", uint(1)).Return(nil)
	ust.On("Verify", mock.Anything).Return(nil)
	us := NewUserService(ust, ver, nil)
	for _, tt := range tests {
		gotErr := us.Verify(tt.giveUID, tt.giveCode)
		if !errors.Is(gotErr, tt.wantErr) {
			t.Errorf("got error: %s, want error: %s", gotErr, tt.wantErr)
		}
	}
}
