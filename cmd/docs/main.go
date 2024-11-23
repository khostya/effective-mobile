package main

import (
	"github.com/khostya/effective-mobile/internal/config"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"log"
	"net"
	"net/http"
	"strconv"
)

func main() {
	cfg := config.MustNewConfig()

	swaggerHandler := httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	)

	addr := net.JoinHostPort("", strconv.Itoa(int(cfg.HTTP.SwaggerPort)))
	if err := http.ListenAndServe(addr, swaggerHandler); err != nil {
		log.Fatalln(err)
	}
}
