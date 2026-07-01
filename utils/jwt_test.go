package utils_test

import (
	"os"
	"testing"
	"time"

	"github.com/Signal-zxh/signalzxh-blog/utils"
	"github.com/golang-jwt/jwt/v5"
)

func setupSecret(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret-key-12345")
}

func TestGenerateToken_Success(t *testing.T) {
	setupSecret(t)

	token, err := utils.GenerateToken(123)
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}

	if token == "" {
		t.Error("GenerateToken() returned empty token")
	}

	tokenParts := len([]rune(token))
	if tokenParts < 100 {
		t.Errorf("GenerateToken() returned unusually short token: %d chars", tokenParts)
	}
}

func TestGenerateToken_DifferentUserIDs(t *testing.T) {
	setupSecret(t)

	token1, err1 := utils.GenerateToken(1)
	token2, err2 := utils.GenerateToken(2)

	if err1 != nil {
		t.Fatalf("GenerateToken(1) error = %v", err1)
	}
	if err2 != nil {
		t.Fatalf("GenerateToken(2) error = %v", err2)
	}

	if token1 == token2 {
		t.Error("GenerateToken() should return different tokens for different user IDs")
	}
}

func TestParseToken_Success(t *testing.T) {
	setupSecret(t)

	userID := 456
	token, err := utils.GenerateToken(userID)
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}

	claims, err := utils.ParseToken(token)
	if err != nil {
		t.Fatalf("ParseToken() error = %v", err)
	}

	if claims == nil {
		t.Fatal("ParseToken() returned nil claims")
	}

	if claims.UserID != userID {
		t.Errorf("ParseToken() UserID = %v, want %v", claims.UserID, userID)
	}

	if claims.ExpiresAt == nil {
		t.Error("ParseToken() ExpiresAt should not be nil")
	}

	if claims.ExpiresAt.Before(time.Now()) {
		t.Error("ParseToken() ExpiresAt should be in the future")
	}
}

func TestParseToken_InvalidToken(t *testing.T) {
	setupSecret(t)

	tests := []struct {
		name      string
		tokenStr  string
		wantError bool
	}{
		{"empty_token", "", true},
		{"invalid_format", "not-a-valid-token", true},
		{"invalid_base64", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.invalid.payload.signature", true},
		{"invalid_signature", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxMjN9.invalid-signature", true},
		{"wrong_secret", generateTokenWithSecret(123, "different-secret"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := utils.ParseToken(tt.tokenStr)
			if (err != nil) != tt.wantError {
				t.Errorf("ParseToken(%q) error = %v, wantError %v", tt.tokenStr, err, tt.wantError)
			}
		})
	}
}

func TestParseToken_ExpiredToken(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret-key-12345")

	claims := utils.Claims{
		UserID: 123,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := []byte(os.Getenv("JWT_SECRET"))
	tokenStr, _ := token.SignedString(secret)

	_, err := utils.ParseToken(tokenStr)
	if err == nil {
		t.Error("ParseToken() should return error for expired token")
	}
}

func TestParseToken_NilToken(t *testing.T) {
	setupSecret(t)

	_, err := utils.ParseToken("")
	if err == nil {
		t.Error("ParseToken() should return error for empty token")
	}
}

func TestClaims_Struct(t *testing.T) {
	claims := &utils.Claims{
		UserID: 123,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}

	if claims.UserID != 123 {
		t.Errorf("Claims.UserID = %v, want 123", claims.UserID)
	}

	if claims.ExpiresAt == nil {
		t.Error("Claims.ExpiresAt should not be nil")
	}
}

func generateTokenWithSecret(userID int, secret string) string {
	claims := utils.Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, _ := token.SignedString([]byte(secret))
	return tokenStr
}

func TestTokenRoundTrip(t *testing.T) {
	setupSecret(t)

	userIDs := []int{1, 100, 999, 12345}

	for _, userID := range userIDs {
		token, err := utils.GenerateToken(userID)
		if err != nil {
			t.Fatalf("GenerateToken(%d) error = %v", userID, err)
		}

		claims, err := utils.ParseToken(token)
		if err != nil {
			t.Fatalf("ParseToken(%d) error = %v", userID, err)
		}

		if claims.UserID != userID {
			t.Errorf("Round trip UserID = %v, want %v", claims.UserID, userID)
		}
	}
}

func TestGenerateToken_ZeroUserID(t *testing.T) {
	setupSecret(t)

	token, err := utils.GenerateToken(0)
	if err != nil {
		t.Fatalf("GenerateToken(0) error = %v", err)
	}

	if token == "" {
		t.Error("GenerateToken(0) returned empty token")
	}

	claims, err := utils.ParseToken(token)
	if err != nil {
		t.Fatalf("ParseToken() error = %v", err)
	}

	if claims.UserID != 0 {
		t.Errorf("ParseToken() UserID = %v, want 0", claims.UserID)
	}
}

func TestParseToken_InvalidClaimsType(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret-key-12345")

	secret := []byte(os.Getenv("JWT_SECRET"))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": 123,
	})

	tokenStr, _ := token.SignedString(secret)

	_, err := utils.ParseToken(tokenStr)
	if err == nil {
		t.Error("ParseToken() should return error for invalid claims type")
	}
}
