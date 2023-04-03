package router

import "github.com/gin-gonic/gin"

func New() (router *gin.Engine) {
	r := gin.Default()
	// r.Use(logging.RequestResponseLogger())

	// gin.SetMode(gin.TestMode)
	r.Use(gin.Recovery())
	r.RedirectTrailingSlash = false
	r.SetTrustedProxies(nil)

	return r
}
