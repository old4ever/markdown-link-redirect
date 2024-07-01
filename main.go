package main

import (
	"fmt"
	"net/http"

	"github.com/aviddiviner/gin-limit"
	"github.com/gin-gonic/contrib/secure"
	"github.com/gin-gonic/gin"
	// "errors"
)

func main() {
	r := gin.Default()
	r.Use(limit.MaxAllowed(1))
	r.Use(secure.Secure(secure.Options{
		// AllowedHosts:          []string{"example.com", "ssl.example.com"},
		// SSLRedirect: true,
		// SSLHost:               "ssl.example.com",
		SSLProxyHeaders:       map[string]string{"X-Forwarded-Proto": "https"},
		STSSeconds:            315360000,
		STSIncludeSubdomains:  true,
		FrameDeny:             true,
		ContentTypeNosniff:    true,
		BrowserXssFilter:      true,
		ContentSecurityPolicy: "default-src 'self'",
	}))
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

	if gin.Mode() == "debug" {
		r.Run(":443")
	} else {
		r.RunTLS(":443", "/etc/letsencrypt/live/magnet.dmytros.dev/fullchain.pem", "/etc/letsencrypt/live/magnet.dmytros.dev/privkey.pem")
	}
}
