package api

import (
	"fmt"
	"gindriver/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func UploadHandler(c *gin.Context) {
	user, ret := c.Get("SessionUser")
	if !ret {
		c.JSON(http.StatusUnauthorized, utils.ErrorWrapper(fmt.Errorf("not authorized")))
		return
	}
	if user != "admin" {
		c.JSON(http.StatusForbidden, utils.ErrorWrapper(fmt.Errorf("forbidden")))
		return
	}

	file, err := c.FormFile("files")
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorWrapper(err))
		return
	}
	err = c.SaveUploadedFile(file, fmt.Sprintf("./public/%s", file.Filename))
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorWrapper(err))
		return
	}
	c.JSON(http.StatusCreated, utils.SuccessWrapper(fmt.Sprintf("file saved at ./public/%s", file.Filename)))
}
