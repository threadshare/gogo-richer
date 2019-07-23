// Copyright (c) 2015 nuboLAB UG (haftungsbeschränkt)
// Use of this source code is governed by the MIT license that can be found
// in the LICENSE file.

// Package jwt provides simple helpers for producing and consuming
// JWT (RFC 7519) that are signed with HMAC SHA-256.
package jwt

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
)

// Registered claim names are defined as constants for convenience.
const (
	Issuer         = "iss"
	Subject        = "sub"
	Audience       = "aud"
	ExpirationTime = "exp"
	NotBefore      = "nbf"
	IssuedAt       = "iat"
	ID             = "jti"
)

// Token contains fields for a JWT.
type Token struct {
	Header    Header
	ClaimSet  ClaimSet
	Signature string
}

// ParseAndVerify verifies a JWT signate and parses the token if ok.
func ParseAndVerify(token, secret string) (Token, bool) {
	if !Verify(token, secret) {
		return Token{}, false
	}
	t, err := Parse(token)
	if err != nil {
		return t, false
	}

	return t, true
}

// Parse parses a JWT without verifying it's signature.
func Parse(token string) (Token, error) {
	var t Token

	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return t, errors.New("jwt: invalid JWT")
	}

	rawHeader, err := base64Decode(parts[0])
	if err != nil {
		return t, err
	}
	if err := json.Unmarshal([]byte(rawHeader), &t.Header); err != nil {
		return t, err
	}

	rawClaims, err := base64Decode(parts[1])
	if err != nil {
		return t, err
	}
	if err := json.Unmarshal([]byte(rawClaims), &t.ClaimSet); err != nil {
		return t, err
	}

	t.Signature = parts[2]

	return t, nil
}

// Header contains information about the token type (always JWT) and the
// algorithm used for signing the token.
type Header struct {
	Type      string `json:"typ"`
	Algorithm string `json:"alg"`
}

// ClaimSet is a map for storing JWT claims.
type ClaimSet map[string]interface{}

// Sign takes a secret and signs the ClaimSet with HMAC SHA-256. It returns
// the base64 encoded byte sequence of the signature or an error in case of
// problems marshalling the ClaimSet to JSON.
func (c ClaimSet) Sign(secret string) (string, error) {
	// XXX hdr directly as base64?
	hdr := `{"typ":"JWT","alg":"HS256"}`
	claims, err := json.Marshal(c)
	if err != nil {
		return "", err
	}
	header := base64Encode([]byte(hdr))
	payload := base64Encode(claims)

	return sign(header+"."+payload, secret), nil
}

// Verify checks the validity of the token and verifies the integrity with
// HMAC SHA-256 and the given secret.
func Verify(token, secret string) bool {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return false
	}

	rawHeader, err := base64Decode(parts[0])
	if err != nil {
		return false
	}
	var hdr Header
	if json.Unmarshal([]byte(rawHeader), &hdr) != nil {
		return false
	}

	if hdr.Type != "JWT" {
		return false
	}

	if hdr.Algorithm != "HS256" {
		return false
	}

	return signature(parts[0]+"."+parts[1], secret) == parts[2]
}

// base64Encode encodes a byte sequence as base64 and strips trailing '='
// from the result.
func base64Encode(b []byte) string {
	return strings.TrimRight(base64.URLEncoding.EncodeToString(b), "=")
}

// base64Decode adds trailing '=' to the input and decodes it from base64.
func base64Decode(s string) ([]byte, error) {
	switch len(s) % 4 {
	case 2:
		s += "=="
	case 3:
		s += "="
	}
	return base64.URLEncoding.DecodeString(s)
}

// sign adds a signature to the message
func sign(message, secret string) string {
	return message + "." + signature(message, secret)
}

// signature returns the signature for a given message as base64 encoded
// byte sequence.
func signature(message, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(message))
	return base64Encode(h.Sum(nil))
}
