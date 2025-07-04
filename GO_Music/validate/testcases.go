package validate

import "time"

type TestStruct struct {
	Birthday time.Time `validate:"birthday_past"`
	Role     string    `validate:"role_valid"`
	Grade    int       `validate:"grade_range"`
	Email    string    `validate:"required,email"`
}

type TestCase struct {
	Name    string
	Input   TestStruct
	WantErr bool
	ErrMsgs map[string]string
}

var Now = time.Now()

var TestCases = []TestCase{
	{
		Name: "valid data",
		Input: TestStruct{
			Birthday: Now.AddDate(-20, 0, 0),
			Role:     "student",
			Grade:    7,
			Email:    "test@example.com",
		},
		WantErr: false,
	},
	{
		Name: "birthday in future",
		Input: TestStruct{
			Birthday: Now.AddDate(1, 0, 0),
			Role:     "admin",
			Grade:    5,
			Email:    "test@example.com",
		},
		WantErr: true,
		ErrMsgs: map[string]string{
			"Birthday": "date must not be in the future",
		},
	},
	{
		Name: "invalid role",
		Input: TestStruct{
			Birthday: Now.AddDate(-30, 0, 0),
			Role:     "manager",
			Grade:    5,
			Email:    "test@example.com",
		},
		WantErr: true,
		ErrMsgs: map[string]string{
			"Role": "role is not valid",
		},
	},
	{
		Name: "grade out of range",
		Input: TestStruct{
			Birthday: Now.AddDate(-30, 0, 0),
			Role:     "teacher",
			Grade:    15,
			Email:    "test@example.com",
		},
		WantErr: true,
		ErrMsgs: map[string]string{
			"Grade": "grade must be between 0 and 100",
		},
	},
	{
		Name: "invalid email",
		Input: TestStruct{
			Birthday: Now.AddDate(-30, 0, 0),
			Role:     "employee",
			Grade:    5,
			Email:    "invalid-email",
		},
		WantErr: true,
		ErrMsgs: map[string]string{
			"Email": "must be a valid email",
		},
	},
	{
		Name: "missing required email",
		Input: TestStruct{
			Birthday: Now.AddDate(-30, 0, 0),
			Role:     "employee",
			Grade:    5,
			Email:    "",
		},
		WantErr: true,
		ErrMsgs: map[string]string{
			"Email": "is required",
		},
	},
}
