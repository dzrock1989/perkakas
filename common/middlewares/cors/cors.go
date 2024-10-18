package cors

import (
	"net/http"
	"strings"

	"github.com/go-chi/cors"
	"github.com/tigapilarmandiri/perkakas"
	"github.com/tigapilarmandiri/perkakas/configs"
)

func Default() func(http.Handler) http.Handler {
	origins := []string{"*"}

	if !perkakas.IsEqual(configs.Config.AllowedOrigins, "*") {
		origins = strings.Split(configs.Config.AllowedOrigins, ",")
	}

	return cors.Handler(cors.Options{
		AllowedOrigins:   origins,
		AllowedMethods:   []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "Host", "Dates"},
		ExposedHeaders:   []string{"Authorization", "Dates"},
		AllowCredentials: false,
		MaxAge:           300,
	})
}
