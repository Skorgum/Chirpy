package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestCheckPasswordHash(t *testing.T) {
	// First, we need to create some hashed passwords for testing
	password1 := "correctPassword123!"
	password2 := "anotherPassword456!"
	hash1, _ := HashPassword(password1)
	hash2, _ := HashPassword(password2)

	tests := []struct {
		name          string
		password      string
		hash          string
		wantErr       bool
		matchPassword bool
	}{
		{
			name:          "Correct password",
			password:      password1,
			hash:          hash1,
			wantErr:       false,
			matchPassword: true,
		},
		{
			name:          "Incorrect password",
			password:      "wrongPassword",
			hash:          hash1,
			wantErr:       false,
			matchPassword: false,
		},
		{
			name:          "Password doesn't match different hash",
			password:      password1,
			hash:          hash2,
			wantErr:       false,
			matchPassword: false,
		},
		{
			name:          "Empty password",
			password:      "",
			hash:          hash1,
			wantErr:       false,
			matchPassword: false,
		},
		{
			name:          "Invalid hash",
			password:      password1,
			hash:          "invalidhash",
			wantErr:       true,
			matchPassword: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			match, err := CheckPasswordHash(tt.password, tt.hash)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckPasswordHash() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && match != tt.matchPassword {
				t.Errorf("CheckPasswordHash() expects %v, got %v", tt.matchPassword, match)
			}
		})
	}
}

func TestValidateJWT(t *testing.T) {
	userID := uuid.New()

	validToken, err := MakeJWT(userID, "secret", time.Hour)
	if err != nil {
		t.Fatalf("Failed to create valid JWT: %v", err)
	}

	expiredToken, err := MakeJWT(userID, "secret", -time.Hour)
	if err != nil {
		t.Fatalf("Failed to create expired JWT: %v", err)
	}

	tests := []struct {
		name        string
		tokenString string
		tokenSecret string
		wantUserID  uuid.UUID
		wantErr     bool
	}{
		{
			name:        "Valid token",
			tokenString: validToken,
			tokenSecret: "secret",
			wantUserID:  userID,
			wantErr:     false,
		},
		{
			name:        "Wrong secret",
			tokenString: validToken,
			tokenSecret: "wrong_secret",
			wantUserID:  uuid.Nil,
			wantErr:     true,
		},
		{
			name:        "Expired token",
			tokenString: expiredToken,
			tokenSecret: "secret",
			wantUserID:  uuid.Nil,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUserID, err := ValidateJWT(tt.tokenString, tt.tokenSecret)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateJWT() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && gotUserID != tt.wantUserID {
				t.Errorf("ValidateJWT() gotUserID = %v, want %v", gotUserID, tt.wantUserID)
			}
		})
	}
}

func TestGetBearerToken(t *testing.T) {
	tests := []struct {
		name      string
		headerVal string
		wantToken string
		wantErr   bool
	}{
		{
			name:      "valid bearer token",
			headerVal: "Bearer abc.def.ghi",
			wantToken: "abc.def.ghi",
			wantErr:   false,
		},
		{
			name:      "missing header",
			headerVal: "",
			wantToken: "",
			wantErr:   true,
		},
		{
			name:      "wrong prefix",
			headerVal: "Token abc.def.ghi",
			wantToken: "",
			wantErr:   true,
		},
		{
			name:      "extra spaces",
			headerVal: "Bearer    abc.def.ghi",
			wantToken: "abc.def.ghi",
			wantErr:   false,
		},
		{
			name:      "no token",
			headerVal: "Bearer",
			wantToken: "",
			wantErr:   true,
		},
		{
			name:      "bearer with only spaces",
			headerVal: "Bearer    ",
			wantToken: "",
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			headers := make(map[string][]string)
			if tt.headerVal != "" {
				headers["Authorization"] = []string{tt.headerVal}
			}
			got, err := GetBearerToken(headers)

			if tt.wantErr && err == nil {
				t.Fatalf("Expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if got != tt.wantToken {
				t.Errorf("expected token %v, got %v", tt.wantToken, got)
			}
		})
	}
}
