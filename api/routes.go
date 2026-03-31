package api

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.Engine) {
	v1 := r.Group("/v1")

	v1.POST("/send", SendEmailHandler)
	v1.GET("/status/:id", GetStatusHandler)
}
