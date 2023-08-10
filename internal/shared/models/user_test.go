package models

import (
	"testing"

	localErrs "github.com/WendelHime/hydroponics-metrics-collector/internal/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestUserBind(t *testing.T) {
	var tests = []struct {
		name      string
		assert    func(t *testing.T, err error)
		givenUser *User
	}{
		{
			name: "bind with success",
			assert: func(t *testing.T, err error) {
				assert.Nil(t, err)
			},
			givenUser: &User{Name: "Test", Email: "test@test.com", Password: "12345678"},
		},
		{
			name: "bind fails if email is invalid",
			assert: func(t *testing.T, err error) {
				if assert.Error(t, err) {
					assert.ErrorIs(t, err, localErrs.BadRequestErr)
				}
			},
			givenUser: &User{Name: "Test", Email: "test", Password: "12345678"},
		},
		{
			name: "bind fails if password doesn't met criteria",
			assert: func(t *testing.T, err error) {
				if assert.Error(t, err) {
					assert.ErrorIs(t, err, localErrs.BadRequestErr)
				}
			},
			givenUser: &User{Name: "Test", Email: "test@test.com", Password: "123456"},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			err := tt.givenUser.Bind(nil)
			tt.assert(t, err)
		})
	}
}
