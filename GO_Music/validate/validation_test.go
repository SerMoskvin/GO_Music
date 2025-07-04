package validate

import (
	"testing"
)

func TestValidateStruct(t *testing.T) {
	for _, tc := range TestCases {
		t.Run(tc.Name, func(t *testing.T) {
			err := ValidateStruct(tc.Input)
			if tc.WantErr {
				if err == nil {
					t.Errorf("Expected error but got nil")
				} else {
					validationErrors, ok := err.(ValidationErrors)
					if !ok {
						t.Errorf("Error should be of type ValidationErrors")
					} else {
						for field, msg := range tc.ErrMsgs {
							if validationErrors[field] != msg {
								t.Errorf("For field '%s', expected message '%s', but got '%s'", field, msg, validationErrors[field])
							}
						}
					}
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				}
			}
		})
	}
}
