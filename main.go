package main

import (
	"fmt"
	"net/http"

	// "github.com/aviddiviner/gin-limit"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"
)

var (
	whitelistedRoutes = map[string]bool{
		"/magnet-redirect": true,
		"/magnet:":         true,
		"/favicon.ico":     true,
		// Add more whitelisted routes here
	}
	blockedIPs = make(map[string]bool)
	logFile    *os.File
	logger     *log.Logger
)

func init() {
	var err error
	logFile, err = os.OpenFile("/var/log/gin_blocked.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	logger = log.New(logFile, "", 0)
}

func WhitelistMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()

		// Check if IP is already blocked
		if blockedIPs[clientIP] {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		// Check if route is whitelisted
		if !whitelistedRoutes[c.Request.URL.Path] {
			// Log the blocked request
			timeStamp := time.Now().Format("2006/01/02 - 15:04:05")
			logger.Printf("[BLOCKED] %s | %s | %s | %s %s\n",
				timeStamp,
				clientIP,
				c.Request.Method,
				c.Request.URL.Path,
				c.Request.URL.RawQuery,
			)

			// Block the IP
			blockedIPs[clientIP] = true

			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		c.Next()
	}
}

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	GIN_MODE := os.Getenv("GIN_MODE")

	r := gin.Default()
	r.Use(WhitelistMiddleware())
	// r.Use(limit.MaxAllowed(1)) // Comment out if using IP whitelist
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
