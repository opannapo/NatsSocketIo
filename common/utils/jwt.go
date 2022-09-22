package utils

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cristalhq/jwt/v3"
	"io"
	"strings"
	"time"
)

var ErrorJwtExpired = errors.New("ERROR_FRAMEWORK_JWT_EXPIRED")

type JwtBuildOption struct {
	JwtSecret              string
	AppName                string
	AppEnv                 string
	JwtTokenExpired        int64
	JwtRefreshTokenExpired int64
}

func JwtBuild(audience []string, option JwtBuildOption) (string, time.Time, string, time.Time, error) {
	key := []byte(option.JwtSecret)
	signer, err := jwt.NewSignerHS(jwt.HS256, key)
	now := time.Now()
	tokenExpired_ := option.JwtTokenExpired
	refreshTokenExpired_ := option.JwtRefreshTokenExpired
	tokenExpired := now.Add(time.Minute * time.Duration(tokenExpired_))
	refreshTokenExpired := now.Add(time.Minute * time.Duration(refreshTokenExpired_))
	if err != nil {
		return "", tokenExpired, "", refreshTokenExpired, err
	}
	claims := &jwt.RegisteredClaims{
		ID:        MD5(fmt.Sprintf("%v", now)),
		Subject:   "at",
		Issuer:    strings.ToLower(fmt.Sprintf("%v_%v", option.AppName, option.AppEnv)),
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(tokenExpired),
	}

	if audience != nil {
		claims.Audience = audience
	}

	builder := jwt.NewBuilder(signer)
	authToken, err := builder.Build(claims)
	if err != nil {
		return "", tokenExpired, "", refreshTokenExpired, err
	}
	claims = &jwt.RegisteredClaims{
		Issuer:    strings.ToLower(fmt.Sprintf("%v_%v", option.AppName, option.AppEnv)),
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(refreshTokenExpired),
		Subject:   "rt",
	}

	refreshToken, err := builder.Build(claims)
	if err != nil {
		return "", tokenExpired, "", refreshTokenExpired, err
	}
	return authToken.String(), tokenExpired, refreshToken.String(), refreshTokenExpired, nil
}

func JwtVerify(text string, jwtSecret string) error {
	key := []byte(jwtSecret)
	verifier, err := jwt.NewVerifierHS(jwt.HS256, key)
	if err != nil {
		return err
	}

	token, err := jwt.ParseString(text)
	if err != nil {
		return err
	}

	err = verifier.Verify(token.Payload(), token.Signature())
	if err != nil {
		return err
	}

	newToken, err := jwt.ParseAndVerifyString(text, verifier)
	if err != nil {
		return err
	}

	// get standard claims
	var newClaims jwt.StandardClaims
	err = json.Unmarshal(newToken.RawClaims(), &newClaims)
	if err != nil {
		return err
	}

	// verify claims as you
	if !newClaims.IsValidAt(time.Now()) {
		return ErrorJwtExpired
	}
	return nil
}

func MD5(plain string) string {
	h := md5.New()
	io.WriteString(h, plain)
	return fmt.Sprintf("%x", h.Sum(nil))
}
