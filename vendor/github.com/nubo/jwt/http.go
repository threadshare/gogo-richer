// Copyright (c) 2015 nuboLAB UG (haftungsbeschränkt)
// Use of this source code is governed by the MIT license that can be found
// in the LICENSE file.
package jwt

import (
	"net/http"
	"strings"
)

func GetRequestHeader(r *http.Request) string {
	raw := r.Header.Get("Authorization")
	if len(raw) < 7 || !strings.EqualFold(raw[:7], "bearer ") {
		return ""
	}
	return raw[7:]
}

func SetRequestHeader(r *http.Request, token string) {
	r.Header.Set("Authorization", "Bearer "+token)
}

func TokenFromRequest(r *http.Request, secret string) (Token, bool) {
	return ParseAndVerify(GetRequestHeader(r), secret)
}
