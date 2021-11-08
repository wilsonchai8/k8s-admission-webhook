package main

import (
	"flag"

	"github.com/gin-gonic/gin"
	"github.com/unrolled/secure"
)

type tlsConfig struct {
	tlsCertFile string
	tlsKeyFile  string
}

func main() {
	tls := tlsConfig{}
	flag.StringVar(&tls.tlsCertFile, "tlsCertFile", "/etc/webhook/certs/cert.pem", "File containing the x509 Certificate for HTTPS.")
	flag.StringVar(&tls.tlsKeyFile, "tlsKeyFile", "/etc/webhook/certs/key.pem", "File containing the x509 private key to --tlsCertFile.")

	flag.Parse()

	router := gin.Default()
	router.Use(TlsHandler())
	router.POST("/mutate", WebhookCallback)

	router.RunTLS(":443", tls.tlsCertFile, tls.tlsKeyFile)
}

func TlsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		secureMiddleware := secure.New(secure.Options{
			SSLRedirect: true,
			SSLHost:     "localhost:443",
		})
		err := secureMiddleware.Process(c.Writer, c.Request)

		if err != nil {
			return
		}

		c.Next()
	}
}
