package validate

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateStruct(t *testing.T) {
	for _, tc := range TestCases {
		t.Run(tc.Name, func(t *testing.T) {
			err := validate.ValidateStruct(tc.Input)
			if tc.WantErr {
				assert.Error(t, err)
				validationErrors, ok := err.(validate.ValidationErrors)
				assert.True(t, ok, "error should be of type ValidationErrors")
				for field, msg := range tc.ErrMsgs {
					assert.Equal(t, msg, validationErrors[field])
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
