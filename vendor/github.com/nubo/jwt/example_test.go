package jwt_test

import (
	"fmt"
	"log"

	"github.com/nubo/jwt"
)

func ExampleClaimSet_Sign() {
	claims := jwt.ClaimSet{
		jwt.Issuer: "example.com",
	}

	token, err := claims.Sign("secret")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(token)
	// Output: eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJleGFtcGxlLmNvbSJ9.6aqSC54aR7dIsuyQgUbcTM4tSkZLcdwqPXzk3OQtOXk
}

func ExampleVerify() {
	token := "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJleGFtcGxlLmNvbSJ9.6aqSC54aR7dIsuyQgUbcTM4tSkZLcdwqPXzk3OQtOXk"

	fmt.Println(jwt.Verify(token, "secret"))
	// Output: true
}

func ExampleParseAndVerify() {
	token := "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJleGFtcGxlLmNvbSJ9.6aqSC54aR7dIsuyQgUbcTM4tSkZLcdwqPXzk3OQtOXk"
	t, ok := jwt.ParseAndVerify(token, "secret")

	fmt.Println(t.Header.Type, t.Header.Algorithm, ok)
	// Output: JWT HS256 true
}

func ExampleParse() {
	token := "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJleGFtcGxlLmNvbSJ9.6aqSC54aR7dIsuyQgUbcTM4tSkZLcdwqPXzk3OQtOXk"
	t, err := jwt.Parse(token)

	fmt.Println(t.Header.Type, t.Header.Algorithm, err)
	// Output: JWT HS256 <nil>
}
