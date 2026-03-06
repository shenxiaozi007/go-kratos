// Package auth 提供 JWT 签发与解析，用于 xsh 后台登录态。
package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims 为 JWT payload，包含用户标识与过期时间。
type Claims struct {
	jwt.RegisteredClaims
	UserID   uint64 `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

// Sign 签发 JWT，expireSeconds 为相对当前的过期秒数。
// 返回 token 字符串与过期时间戳（秒），供前端或 LoginReply 使用。
func Sign(secret string, userID uint64, username, role string, expireSeconds int32) (token string, expAt int64, err error) {
	if expireSeconds <= 0 {
		expireSeconds = 86400
	}
	now := time.Now()
	exp := now.Add(time.Duration(expireSeconds) * time.Second)
	expAt = exp.Unix()
	claims := &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
			IssuedAt:  jwt.NewNumericDate(now),
		},
		UserID:   userID,
		Username: username,
		Role:     role,
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err = t.SignedString([]byte(secret))
	if err != nil {
		return "", 0, fmt.Errorf("jwt sign: %w", err)
	}
	return token, expAt, nil
}

// Parse 解析并校验 JWT，返回 Claims；过期或签名无效返回错误。
func Parse(secret, tokenString string) (*Claims, error) {
	t, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("jwt parse: %w", err)
	}
	claims, ok := t.Claims.(*Claims)
	if !ok || !t.Valid {
		return nil, fmt.Errorf("jwt invalid")
	}
	return claims, nil
}
