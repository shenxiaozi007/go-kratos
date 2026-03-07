package auth

import (
	"testing"
	"time"
)

func TestSign(t *testing.T) {
	secret := "test-secret"
	userID := uint64(1)
	username := "alice"
	role := "admin"
	expireSeconds := int32(3600)

	t.Run("valid_sign_returns_token_and_exp", func(t *testing.T) {
		token, expAt, err := Sign(secret, userID, username, role, expireSeconds)
		if err != nil {
			t.Fatalf("Sign: %v", err)
		}
		if token == "" {
			t.Error("token is empty")
		}
		now := time.Now().Unix()
		if expAt <= now {
			t.Errorf("expAt %d should be after now %d", expAt, now)
		}
		if expAt > now+int64(expireSeconds)+2 {
			t.Errorf("expAt %d should be about now+%d", expAt, expireSeconds)
		}
	})

	t.Run("zero_expire_uses_default", func(t *testing.T) {
		token, expAt, err := Sign(secret, userID, username, role, 0)
		if err != nil {
			t.Fatalf("Sign(0): %v", err)
		}
		if token == "" {
			t.Error("token is empty")
		}
		// 默认 86400 秒
		now := time.Now().Unix()
		if expAt <= now+3600 {
			t.Error("with default expire, expAt should be ~24h later")
		}
		_ = expAt
	})
}

func TestParse(t *testing.T) {
	secret := "parse-test-secret"
	token, _, err := Sign(secret, 1, "bob", "operator", 3600)
	if err != nil {
		t.Fatalf("Sign for Parse test: %v", err)
	}

	tests := []struct {
		name    string
		secret  string
		token   string
		wantErr bool
		check   func(t *testing.T, c *Claims)
	}{
		{
			name:    "valid_token",
			secret:  secret,
			token:   token,
			wantErr: false,
			check: func(t *testing.T, c *Claims) {
				if c.UserID != 1 {
					t.Errorf("UserID = %d, want 1", c.UserID)
				}
				if c.Username != "bob" {
					t.Errorf("Username = %s, want bob", c.Username)
				}
				if c.Role != "operator" {
					t.Errorf("Role = %s, want operator", c.Role)
				}
			},
		},
		{
			name:    "wrong_secret",
			secret:  "other-secret",
			token:   token,
			wantErr: true,
			check:   nil,
		},
		{
			name:    "empty_token",
			secret:  secret,
			token:   "",
			wantErr: true,
			check:   nil,
		},
		{
			name:    "malformed_token",
			secret:  secret,
			token:   "not.a.valid.jwt",
			wantErr: true,
			check:   nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := Parse(tt.secret, tt.token)
			if tt.wantErr {
				if err == nil {
					t.Error("Parse: expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("Parse: %v", err)
			}
			if claims == nil {
				t.Fatal("claims is nil")
			}
			if tt.check != nil {
				tt.check(t, claims)
			}
		})
	}
}

func TestSignParseRoundTrip(t *testing.T) {
	secret := "roundtrip-secret"
	userID := uint64(42)
	username := "roundtrip-user"
	role := "admin"

	token, _, err := Sign(secret, userID, username, role, 86400)
	if err != nil {
		t.Fatalf("Sign: %v", err)
	}
	claims, err := Parse(secret, token)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if claims.UserID != userID || claims.Username != username || claims.Role != role {
		t.Errorf("roundtrip: got UserID=%d Username=%s Role=%s", claims.UserID, claims.Username, claims.Role)
	}
}
