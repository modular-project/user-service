package middleware

import (
	"fmt"
	"net/http"
	"testing"
	"users-service/http/handler"
	"users-service/pkg"
)

func TestFindError(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr error
	}{
		{
			name:    "BadRequest",
			args:    args{fmt.Errorf("error, %w", fmt.Errorf("error 2, %w", handler.ErrBindData))},
			want:    http.StatusBadRequest,
			wantErr: handler.ErrBindData,
		}, {
			name: "Internal",
			args: args{
				err: fmt.Errorf("a undefinied error"),
			},
			want:    http.StatusInternalServerError,
			wantErr: nil,
		}, {
			name:    "Unauthorized",
			args:    args{fmt.Errorf("error, %w", fmt.Errorf("error 2, %w", pkg.UnauthorizedErr("user is unauthorized")))},
			want:    http.StatusUnauthorized,
			wantErr: pkg.UnauthorizedErr("user is unauthorized"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := findError(tt.args.err)
			if err != tt.wantErr {
				t.Errorf("FindError() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("FindError() = %v, want %v", got, tt.want)
			}
		})
	}
}
