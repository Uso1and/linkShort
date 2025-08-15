package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func IndexPageHandler(c *gin.Context) {

	c.HTML(http.StatusOK, "index.html", nil)
}

func RegisterPageHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "register.html", nil)
}

func LoginPageHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}
