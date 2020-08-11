package api

import (
	"fmt"
	"gindriver/models"
	"gindriver/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func BeginRegistration(c *gin.Context) {
	var (
		beginRegRequest models.BeginRegRequest
		err             error
	)

	err = c.BindJSON(&beginRegRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorWrapper(fmt.Errorf("bad json body")))
		return
	}

	username := beginRegRequest.Username
	if !utils.FilterUsername(username) {
		c.JSON(http.StatusBadRequest, utils.ErrorWrapper(fmt.Errorf("bad username")))
		return
	}

	ret, err := models.GetValidUserByName(username)
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
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorWrapper(err))
		return
	}

	err = utils.WebAuthnSessionStore.SaveWebauthnSession("reg", sessionData, c.Request, c.Writer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorWrapper(err))
		return
	}

	c.JSON(http.StatusCreated, options)
}

func FinishRegistration(c *gin.Context) {
	username := c.Param("name")
	if !utils.FilterUsername(username) {
		c.JSON(http.StatusBadRequest, utils.ErrorWrapper(fmt.Errorf("bad username")))
		return
	}

	user, err := models.GetUserByName(username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorWrapper(err))
		return
	}

	if user.Name == "" {
		c.JSON(http.StatusBadRequest, utils.ErrorWrapper(fmt.Errorf("user not exist")))
		return
	}

	sessionData, err := utils.WebAuthnSessionStore.GetWebauthnSession("reg", c.Request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorWrapper(err))
		return
	}

	credential, err := utils.WebAuthn.FinishRegistration(user, sessionData, c.Request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorWrapper(err))
		return
	}

	user.Valid = true
	user.CredentialId = utils.Btoa(credential.ID)
	user.CredentialAttestationType = credential.AttestationType
	user.AuthenticatorAAGUID = utils.Btoa(credential.Authenticator.AAGUID)
	user.AuthenticatorSignCount = credential.Authenticator.SignCount

	err = utils.SaveFile(fmt.Sprintf("./public/pubkeys/%s.pub", username), credential.PublicKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorWrapper(err))
		return
	}

	utils.Database.Save(&user)
	c.JSON(http.StatusOK, utils.SuccessWrapper("Registration Success"))
}

func BeginLogin(c *gin.Context) {
	username := c.Param("name")
	if !utils.FilterUsername(username) {
		c.JSON(http.StatusBadRequest, utils.ErrorWrapper(fmt.Errorf("bad username")))
		return
	}

	user, err := models.GetValidUserByName(username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorWrapper(err))
		return
	}
	if user.Name == "" {
		c.JSON(http.StatusBadRequest, utils.ErrorWrapper(fmt.Errorf("user not exist")))
		return
	}

	queryUser := user.ConstructWebAuthNUser()

	options, sessionData, err := utils.WebAuthn.BeginLogin(queryUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorWrapper(err))
		return
	}

	if err = utils.WebAuthnSessionStore.
		SaveWebauthnSession("auth", sessionData, c.Request, c.Writer); err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorWrapper(err))
		return
	}

	c.JSON(http.StatusOK, options)
}

func FinishLogin(c *gin.Context) {
	username := c.Param("name")
	if !utils.FilterUsername(username) {
		c.JSON(http.StatusBadRequest, utils.ErrorWrapper(fmt.Errorf("bad username")))
		return
	}

	user, err := models.GetValidUserByName(username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorWrapper(err))
		return
	}
	if user.Name == "" {
		c.JSON(http.StatusBadRequest, utils.ErrorWrapper(fmt.Errorf("user not exist")))
		return
	}

	sessionData, err := utils.WebAuthnSessionStore.GetWebauthnSession("auth", c.Request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorWrapper(err))
		return
	}

	queryUser := user.ConstructWebAuthNUser()

	_, err = utils.WebAuthn.FinishLogin(queryUser, sessionData, c.Request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorWrapper(err))
		return
	}

	jwtToken, err := models.GenerateJWTToken(user.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorWrapper(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login success",
		"token":   jwtToken,
	})
}
