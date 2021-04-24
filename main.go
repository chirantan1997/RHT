package main

import (
	router "Newton/routers"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/rs/cors"
)

func main() {
	r := router.Router()
	http.Handle("/", r)
	fmt.Println("Starting Server.....")
	fmt.Println("Listening on Port 8080......")

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"*"},
	})

	handler := c.Handler(r)
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), handler))

}
