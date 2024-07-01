package main

import (
	"fmt"
	"net/http"

	"github.com/aviddiviner/gin-limit"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	GIN_MODE := os.Getenv("GIN_MODE")

	r := gin.Default()
	r.Use(limit.MaxAllowed(1))
	r.SetTrustedProxies(nil)

	r.GET("/magnet-redirect", func(ctx *gin.Context) {
		hash := ctx.Query("hash")
		magnet := fmt.Sprintf("magnet:?xt=urn:btih:%s", hash)
		if hash == "" {
			ctx.String(http.StatusBadRequest, "hash query parameter is missing")
		} else {
			ctx.Redirect(http.StatusFound, magnet)
		}
	})

	if GIN_MODE == "debug" {
		r.Run(":443")
	} else {
		r.RunTLS(":443", "/etc/letsencrypt/live/magnet.dmytros.dev/fullchain.pem", "/etc/letsencrypt/live/magnet.dmytros.dev/privkey.pem")
	}
}
