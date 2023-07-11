package main

import (
	"log"
	"posts_ms/internal/database"
	"posts_ms/internal/posts"

	"github.com/gofiber/fiber/v2"
)

func main() {
	database.ConnectDB()

	app := fiber.New()

	app.Post("/new-post", posts.NewPost)

	log.Fatal(app.Listen("127.0.0.1:5000"))
}
