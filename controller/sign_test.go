package controller

import (
	"errors"
	"fmt"
	"net/http"
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
		give     string
		WantID   uint
		WantCode int
	}{
		{
			give:   "OK",
			WantID: 1,
		}, {
			give:     "Bad",
			WantID:   0,
			WantCode: http.StatusBadRequest,
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
		gotID, err := ss.validateRefreshToken(&tt.give)
		gotCode, _ := pkg.FindError(err)
		assert.Equal(t, gotID, tt.WantID)
		assert.Equal(t, tt.WantCode, gotCode)
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
	errNotFound := errors.New("not found")
	tests := []struct {
		giveUID  uint
		giveCode string
		wantCode int
	}{
		{
			giveUID:  1,
			giveCode: "CODEBAD",
			wantCode: http.StatusBadRequest,
		}, {
			giveUID:  1,
			giveCode: "CODEOK",
		}, {
			giveUID:  2,
			giveCode: "CODEOK",
			wantCode: http.StatusBadRequest,
		}, {
			giveUID:  3,
			giveCode: "EXPIREDCODE",
			wantCode: http.StatusBadRequest,
		},
	}
	ust := mocks.NewUserStorager(t)
	ver := mocks.NewVerificationStorager(t)
	ver.On("Find", uint(1)).Return(model.Verification{Code: "CODEOK", ExpiresAt: time.Now().Add(1 * time.Minute)}, nil)
	ver.On("Find", uint(2)).Return(model.Verification{}, errNotFound)
	ver.On("Find", uint(3)).Return(model.Verification{Code: "EXPIREDCODE", ExpiresAt: time.Now().Add(-1 * time.Minute)}, nil)
	ver.On("Delete", uint(1)).Return(nil)
	ust.On("Verify", mock.Anything).Return(nil)
	us := NewUserService(ust, ver, nil)
	for _, tt := range tests {
		err := us.Verify(tt.giveUID, tt.giveCode)
		gotCode, _ := pkg.FindError(err)
		assert.Equal(t, tt.wantCode, gotCode)
	}
}

func TestIsNotFoundErr(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "true",
			args: args{
				err: fmt.Errorf("this is a test: %w", pkg.ErrNoRowsAffected),
			},
			want: true,
		}, {
			name: "false",
			args: args{
				err: fmt.Errorf("this is a test:"),
			},
			want: false,
		}, {
			name: "nil",
			args: args{nil},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := errors.Is(tt.args.err, pkg.ErrNoRowsAffected); got != tt.want {
				t.Errorf("ispkg.NotFoundErr() = %v, want %v", got, tt.want)
			}
		})
	}
}
