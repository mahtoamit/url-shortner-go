package routes

import (
    "github.com/gofiber/fiber/v2"
    "github.com/mahtoamit/urlshortnerservice/api/redis-db"
    "os"
    "time"
    "github.com/go-redis/redis/v8"
    "math/rand"
    "strings"
    "net/url"
)

func ShortenURL(c *fiber.Ctx) error {
    type request struct {
        URL string `json:"url"`
    }
    var body request
    if err := c.BodyParser(&body); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse JSON"})
    }
    if !isValidURL(body.URL) {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid URL format"})
    }
    id := generateRandomString()

    expiry, err := time.ParseDuration(os.Getenv("URL_EXPIRY"))
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "invalid expiry duration"})
    }

    if err := redisdb.Rdb.Set(redisdb.Ctx, id, body.URL, expiry).Err(); err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not store URL"})
    }

    baseURL := os.Getenv("BASE_URL")
    if baseURL == "" {
        baseURL = "http://localhost:3000" // default value if BASE_URL is not set
    }

    return c.Status(fiber.StatusOK).JSON(fiber.Map{"shortened_url": baseURL + "/" + id})
}

func ResolveURL(c *fiber.Ctx) error {
    id := c.Params("id")

    url, err := redisdb.Rdb.Get(redisdb.Ctx, id).Result()
    if err == redis.Nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "URL not found"})
    } else if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not retrieve URL"})
    }

    return c.Redirect(url, fiber.StatusMovedPermanently)
}

func generateRandomString() string {
    const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
    var b strings.Builder
    for i := 0; i < 5; i++ {
        b.WriteByte(charset[rand.Intn(len(charset))])
    }
    for i := 0; i < 2; i++ {
        digit := rand.Intn(10)
        b.WriteByte(byte('0' + digit))
    }
    return b.String()
}

func isValidURL(u string) bool {
    _, err := url.ParseRequestURI(u)
    if err != nil {
        return false
    }
    return true
}
