package api

import (
	"fmt"
	"gindriver/models"
	"gindriver/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func BeginRegistration(c *gin.Context) {
	username := c.PostForm("username")
	var err error

	if !utils.FilterUsername(username) {
		c.JSON(http.StatusBadRequest, utils.ErrorWrapper(fmt.Errorf("bad username")))
		return
	}

	ret, err := models.GetUserByName(username)
	if err != nil && err.Error() != "record not found" {
		c.JSON(http.StatusInternalServerError, utils.ErrorWrapper(err))
		return
	}
	if ret.Name != "" {
		c.JSON(http.StatusBadRequest, utils.ErrorWrapper(fmt.Errorf("user exist")))
		return
	}

	newUser := models.NewUser(username, username)
	if newUser, err = newUser.Insert(); err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorWrapper(err))
		return
	}

	options, sessionData, err := utils.WebAuthn.BeginRegistration(newUser)

	fmt.Printf("Session data: %s", sessionData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorWrapper(err))
		return
	}

	c.JSON(http.StatusOK, options)
}
