// Copyright (c) 2015 nuboLAB UG (haftungsbeschränkt)
// Use of this source code is governed by the MIT license that can be found
// in the LICENSE file.

package jwt_test

import (
	"strings"

	"testing"

	"github.com/nubo/jwt"
)

func TestClaimSet(t *testing.T) {
	claims := jwt.ClaimSet{
		jwt.Issuer: "example.com",
	}
	s, err := claims.Sign("secret")
	if err != nil {
		t.Errorf("Error signing claims %q, %q", claims, err)
	}

	parts := strings.Split(s, ".")
	if len(parts) != 3 {
		t.Errorf("Invalid token %q", s)
	}

	if parts[0] != "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9" {
		t.Errorf("Invalid header: %q", parts[0])
	}

	if parts[1] != "eyJpc3MiOiJleGFtcGxlLmNvbSJ9" {
		t.Errorf("Invalid claims: %q", parts[1])
	}

	if parts[2] != "6aqSC54aR7dIsuyQgUbcTM4tSkZLcdwqPXzk3OQtOXk" {
		t.Errorf("Invalid signature: %q", parts[2])
	}
}

func TestTokenVerification(t *testing.T) {
	token := "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJleGFtcGxlLmNvbSJ9.6aqSC54aR7dIsuyQgUbcTM4tSkZLcdwqPXzk3OQtOXk"
	if !jwt.Verify(token, "secret") {
		t.Error("Verify failed")
	}
}

func TestTokenParse(t *testing.T) {
	data := "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJleGFtcGxlLmNvbSJ9.6aqSC54aR7dIsuyQgUbcTM4tSkZLcdwqPXzk3OQtOXk"
	var token jwt.Token
	token, err := jwt.Parse(data)
	if err != nil {
		t.Error(err)
	}
	if token.Header.Type != "JWT" {
		t.Errorf("Invalid token type: %q", token.Header.Type)
	}
	if token.Header.Algorithm != "HS256" {
		t.Errorf("Invalid token algorithm: %q", token.Header.Algorithm)
	}
	issuer, ok := token.ClaimSet[jwt.Issuer]
	if !ok {
		t.Errorf("Issues missing in claim set: %q", token.ClaimSet)
	}
	if issuer != "example.com" {
		t.Errorf("Invalid issues: %q", token.ClaimSet[jwt.Issuer])
	}
	if token.Signature != "6aqSC54aR7dIsuyQgUbcTM4tSkZLcdwqPXzk3OQtOXk" {
		t.Errorf("Signature not ok: %q", token.Signature)
	}
}

func TestTokenParseAndVerify(t *testing.T) {
	data := "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJleGFtcGxlLmNvbSJ9.6aqSC54aR7dIsuyQgUbcTM4tSkZLcdwqPXzk3OQtOXk"
	var token jwt.Token
	token, ok := jwt.ParseAndVerify(data, "secret")
	if !ok {
		t.Error("ParseAndVerify failed")
	}
	if token.Header.Type != "JWT" {
		t.Errorf("Invalid token type: %q", token.Header.Type)
	}
	if token.Header.Algorithm != "HS256" {
		t.Errorf("Invalid token algorithm: %q", token.Header.Algorithm)
	}
	issuer, ok := token.ClaimSet[jwt.Issuer]
	if !ok {
		t.Errorf("Issues missing in claim set: %q", token.ClaimSet)
	}
	if issuer != "example.com" {
		t.Errorf("Invalid issues: %q", token.ClaimSet[jwt.Issuer])
	}
	if token.Signature != "6aqSC54aR7dIsuyQgUbcTM4tSkZLcdwqPXzk3OQtOXk" {
		t.Errorf("Signature not ok: %q", token.Signature)
	}
}
