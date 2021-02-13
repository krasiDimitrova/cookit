package main

import (
	"github.com/krasimiraMilkova/cookit/internal/routes"
	"github.com/krasimiraMilkova/cookit/internal/users/auth"
	"github.com/rs/cors"
	"log"
	"net/http"
)

func main() {
	router := routes.Handlers()

	http.Handle("/", router)
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedHeaders:   []string{auth.TokenName},
	})

	handler := c.Handler(router)

	err := http.ListenAndServe(":8080", handler)
	if err != nil {
		log.Fatal(err)
	}
}
