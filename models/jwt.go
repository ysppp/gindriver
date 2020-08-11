package models

import (
	"fmt"
	"gindriver/utils"
	"github.com/dgrijalva/jwt-go"
)

type JWTClaims struct {
	Name string `json:"name"`
	jwt.StandardClaims
}

func GenerateJWTToken(name string) (token string, err error) {
	unsignedToken := jwt.NewWithClaims(jwt.SigningMethodHS512, JWTClaims{
		Name: name,
	})
	return unsignedToken.SignedString(utils.ReadFile(fmt.Sprintf("./public/pubkeys/%s.pub", name)))
}
