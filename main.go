package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	// "errors"
)

func main() {
	r := gin.Default()

	// r.GET("/ping", func(ctx *gin.Context) {
	// 	ctx.JSON(http.StatusOK, gin.H{
	// 		"message": "pong",
	// 	})
	// })

	// r.GET("/redirect", func(ctx *gin.Context) {
	// 	ctx.Redirect(http.StatusFound, "/ping")
	// })

	// r.GET("/query", func(ctx *gin.Context) {
	// 	hash := ctx.Query("hash")
	// 	if hash == "" {
	// 		ctx.String(http.StatusBadRequest, "hash query parameter is missing")
	// 	} else {
	// 		ctx.String(http.StatusOK, "Hash is: %s", hash)
	// 	}
	// })

	// r.GET("/", func(ctx *gin.Context) {
	// 	ctx.String(http.StatusNoContent, "send hash to /magnet-redirect?hash=")
	// })

	r.GET("/magnet-redirect", func(ctx *gin.Context) {
		hash := ctx.Query("hash")
		magnet := fmt.Sprintf("magnet:?xt=urn:btih:%s", hash)
		if hash == "" {
			ctx.String(http.StatusBadRequest, "hash query parameter is missing")
		} else {
			ctx.Redirect(http.StatusFound, magnet)
		}
	})

	// r.Run(":80")
	r.Run(":8080")
}
