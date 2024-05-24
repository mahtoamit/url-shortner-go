package main

import (
    "github.com/gofiber/fiber/v2"
    "github.com/joho/godotenv"
    "github.com/mahtoamit/urlshortnerservice/api/redis-db"
    "github.com/mahtoamit/urlshortnerservice/api/routes"
    "log"
)

func main() {
    if err := godotenv.Load(); err != nil {
        log.Fatalf("Error loading .env file: %v", err)
    }
    redisdb.InitRedis()
    app := fiber.New()
    // Serve frontend files
    app.Static("/", "./frontend")

    // API endpoints
    app.Post("/shorten", routes.ShortenURL)
    app.Get("/:id", routes.ResolveURL)
    log.Fatal(app.Listen(":3000"))
}
