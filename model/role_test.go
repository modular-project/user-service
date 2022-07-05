package model_test

import (
	"fmt"
	"testing"
	"users-service/model"

	"github.com/stretchr/testify/assert"
)

func TestIsGreater(t *testing.T) {
	tests := []struct {
		give       model.RoleID
		giveTarget model.RoleID
		want       bool
	}{
		{
			give:       model.OWNER,
			giveTarget: model.ADMIN,
			want:       true,
		}, {
			give:       model.OWNER,
			giveTarget: model.OWNER,
			want:       false,
		}, {
			give:       model.USER,
			giveTarget: model.WAITER,
			want:       false,
		}, {
			give:       model.WAITER,
			giveTarget: model.USER,
			want:       true,
		}, {
			give:       model.OWNER,
			giveTarget: model.USER,
			want:       true,
		},
	}
	for _, tt := range tests {
		assert.Equal(t, tt.want, tt.give.IsGreater(tt.giveTarget), fmt.Sprintf("give: %d, giveTarget: %d", tt.give, tt.giveTarget))
	}
}
