package db_test

import (
	"reflect"
	"testing"
	"time"

	"GO_Music/db"
)

type TestUser struct {
	ID        int       `db:"user_id"`
	Name      string    `db:"user_name"`
	Email     string    `db:"email_address"`
	CreatedAt time.Time `db:"created_at"`
	IsActive  bool      `db:"is_active"`
	secret    string    // неэкспортируемое поле (должно игнорироваться)
}

type TestProduct struct {
	ProductID   int
	ProductName string
	Price       float64
}

func TestToSnakeCase(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Simple case", "ID", "id"},
		{"Two words", "UserName", "user_name"},
		{"Three words", "EmailAddress", "email_address"},
		{"With acronym", "HTMLParser", "html_parser"},
		{"Already snake", "already_snake", "already_snake"},
		{"Mixed case", "TestABC123", "test_abc123"},
		{"Empty string", "", ""},
		{"Another mix", "Model3D", "model3_d"},
		{"Mix register", "userID", "user_id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := db.ToSnakeCase(tt.input)
			if result != tt.expected {
				t.Errorf("toSnakeCase(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestStructToMap(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name     string
		input    interface{}
		expected map[string]interface{}
		wantErr  bool
	}{
		{
			name: "Simple struct with tags",
			input: TestUser{
				ID:        1,
				Name:      "John",
				Email:     "john@example.com",
				CreatedAt: now,
				IsActive:  true,
				secret:    "hidden",
			},
			expected: map[string]interface{}{
				"user_id":       1,
				"user_name":     "John",
				"email_address": "john@example.com",
				"created_at":    now,
				"is_active":     true,
			},
			wantErr: false,
		},
		{
			name: "Pointer to struct",
			input: &TestUser{
				ID:    2,
				Name:  "Alice",
				Email: "alice@example.com",
			},
			expected: map[string]interface{}{
				"user_id":       2,
				"user_name":     "Alice",
				"email_address": "alice@example.com",
				"created_at":    time.Time{},
				"is_active":     false,
			},
			wantErr: false,
		},
		{
			name: "Struct without tags",
			input: TestProduct{
				ProductID:   101,
				ProductName: "Laptop",
				Price:       999.99,
			},
			expected: map[string]interface{}{
				"product_id":   101,
				"product_name": "Laptop",
				"price":        999.99,
			},
			wantErr: false,
		},
		{
			name:     "Non-struct input",
			input:    "not a struct",
			expected: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := db.StructToMap(tt.input)

			if (err != nil) != tt.wantErr {
				t.Errorf("StructToMap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("StructToMap() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestMapToStruct(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name     string
		input    map[string]interface{}
		target   interface{}
		expected interface{}
		wantErr  bool
	}{
		{
			name: "Simple case with tags",
			input: map[string]interface{}{
				"user_id":       1,
				"user_name":     "John",
				"email_address": "john@example.com",
				"created_at":    now,
				"is_active":     true,
			},
			target: &TestUser{},
			expected: &TestUser{
				ID:        1,
				Name:      "John",
				Email:     "john@example.com",
				CreatedAt: now,
				IsActive:  true,
			},
			wantErr: false,
		},
		{
			name: "With snake_case conversion",
			input: map[string]interface{}{
				"product_id":   101,
				"product_name": "Laptop",
				"price":        999.99,
			},
			target: &TestProduct{},
			expected: &TestProduct{
				ProductID:   101,
				ProductName: "Laptop",
				Price:       999.99,
			},
			wantErr: false,
		},
		{
			name: "Non-pointer target",
			input: map[string]interface{}{
				"user_id": 1,
			},
			target:   TestUser{},
			expected: TestUser{},
			wantErr:  true,
		},
		{
			name: "Nil target",
			input: map[string]interface{}{
				"user_id": 1,
			},
			target:   nil,
			expected: nil,
			wantErr:  true,
		},
		{
			name: "Type conversion",
			input: map[string]interface{}{
				"user_id":   int64(1),       // different type
				"user_name": []byte("John"), // []byte to string
				"is_active": int64(1),       // int to bool
			},
			target: &TestUser{},
			expected: &TestUser{
				ID:       1,
				Name:     "John",
				IsActive: true,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := db.MapToStruct(tt.input, tt.target)

			if (err != nil) != tt.wantErr {
				t.Errorf("MapToStruct() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && !reflect.DeepEqual(tt.target, tt.expected) {
				t.Errorf("MapToStruct() = %+v, want %+v", tt.target, tt.expected)
			}
		})
	}
}
