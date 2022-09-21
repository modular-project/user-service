package middleware

import (
	"fmt"
	"net/http"
	"testing"
	"users-service/pkg"

	"github.com/stretchr/testify/assert"
)

func TestFindError(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantMSG string
	}{
		{
			name:    "BadRequest",
			args:    args{fmt.Errorf("error, %w", fmt.Errorf("error 2, %w", &pkg.AppError{MSG: "BadRequest", Code: http.StatusBadRequest}))},
			want:    http.StatusBadRequest,
			wantMSG: "BadRequest",
		}, {
			name: "Internal",
			args: args{
				err: fmt.Errorf("a undefinied error"),
			},
			want: http.StatusInternalServerError,
		}, {
			name:    "Unauthorized",
			args:    args{fmt.Errorf("error, %w", fmt.Errorf("error 2, %w", pkg.NewAppError("Unauthorized", nil, http.StatusUnauthorized)))},
			want:    http.StatusUnauthorized,
			wantMSG: "Unauthorized",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotMSG := pkg.FindError(tt.args.err)
			assert.Equal(t, tt.want, got, tt.name)
			assert.Equal(t, tt.wantMSG, gotMSG, tt.name)
		})
	}
}
