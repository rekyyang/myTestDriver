package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	svc := gin.New()
	svc.GET("/v1/:prefix/*suffix", func(context *gin.Context) {
		context.JSON(http.StatusOK, struct {
			V string
			W string
		}{
			V: context.Param("prefix"),
			W: context.Param("suffix"),
		})
	})
	svc.Run("")
}
