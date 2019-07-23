# JWT - Helpers for Go [![Circle CI](https://circleci.com/gh/nubo/jwt/tree/develop.svg?style=svg)](https://circleci.com/gh/nubo/jwt/tree/develop) [![GoDoc](https://godoc.org/github.com/nubo/jwt?status.svg)](https://godoc.org/github.com/nubo/jwt)

This package implements helpers for handling HMAC SHA-256 signed JSON Web Tokens ([RFC 7519](https://tools.ietf.org/html/rfc7519)) in Go.

## Usage

### Producing a Token

```go
package main

import (
  "fmt"
  "log"

  "github.com/nubo/jwt"
)

func main() {
  claims := jwt.ClaimSet{
    jwt.Issuer:   "example.com",
    jwt.Audience: "example.com",
    "lorem":      "ipsum",
  }
  token, err := claims.Sign("secret")
  if err != nil {
    log.Fatal(err)
  }
  fmt.Println(token)
}
```

### Consuming a Token

```go
package main

import (
  "fmt"
  "log"

  "github.com/nubo/jwt"
)

func main() {
	rawToken := "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJhdWQiOiJleGFtcGxlLmNvbSIsImlzcyI6ImV4YW1wbGUuY29tIiwibG9yZW0iOiJpcHN1bSJ9.VhJwcvoGPhr_sY_YG6-rMNwU0YnpDSGw7jlArsnj8eA"

	token, ok := jwt.ParseAndVerify(rawToken, "secret")
	if !ok {
		log.Fatal("Invalid token")
	}
	fmt.Println("Type", token.Header.Type)
	fmt.Println("Algorithm", token.Header.Algorithm)
	fmt.Println("Claim Set", token.ClaimSet)
	fmt.Println("Signature", token.Signature)
}
```

## Features

There is currently no plan to implement other signing algorithms than HMAC
SHA-256.

- [x] sign JWT with HMAC SHA-256
- [ ] encrypt JWT (JWE)
- [x] parse a raw JWT to a Go struct
- [x] verify a raw JWT
- [x] parse and verify a raw JWT to a Go struct
- [x] claim set
