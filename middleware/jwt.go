package middleware

import (
	"encoding/json"
	"fmt"
	"gindriver/models"
	"gindriver/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func LoginRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Testing admin token
		// eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.
		// eyJuYW1lIjoiYWRtaW4ifQ.
		// nNvEZFA3VmUKjE8cSZIgXZj7ETRkyxNLDtAmBzaC77UNzEiWg31zrgmq44CqOmxLuyvfLyKdPFKBhxmimSbbDQ

		authHeader := c.Request.Header.Get("Authorization")
		if len(authHeader) == 0 {
			c.JSON(http.StatusUnauthorized, utils.ErrorWrapper(fmt.Errorf("unauthorized")))
			c.Abort()
			return
		}

		tmp := strings.Split(authHeader, " ")
		if len(tmp) != 2 {
			c.JSON(http.StatusBadRequest, utils.ErrorWrapper(fmt.Errorf("bad authentication token")))
			c.Abort()
			return
		}
		jwtToken := tmp[1]
		payload := strings.Split(jwtToken, ".")
		if len(payload) != 3 {
			c.JSON(http.StatusBadRequest, utils.ErrorWrapper(fmt.Errorf("bad authentication token")))
			c.Abort()
			return
		}
		jsonPayload := &models.JWTClaims{}
		err := json.Unmarshal(utils.AtobURLSafe(payload[1]), &jsonPayload)
		if err != nil {
			c.JSON(http.StatusBadRequest, utils.ErrorWrapper(fmt.Errorf("bad authentication token")))
			c.Abort()
			return
		}

		token, err := jwt.ParseWithClaims(jwtToken, &models.JWTClaims{},
			func(token *jwt.Token) (interface{}, error) {
				return utils.ReadFile(fmt.Sprintf("./public/pubkeys/%s.pub", jsonPayload.Name)), nil
			})
		if err != nil {
			c.JSON(http.StatusUnauthorized, utils.ErrorWrapper(fmt.Errorf("bad credential")))
			c.Abort()
			return
		}

		if _, err := token.Claims.(*models.JWTClaims); err && token.Valid {
			c.Set("SessionUser", jsonPayload.Name)
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, utils.ErrorWrapper(fmt.Errorf("bad credential")))
			c.Abort()
			return
		}
	}
}
