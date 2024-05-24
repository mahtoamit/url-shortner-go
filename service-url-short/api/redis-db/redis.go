package redisdb

import (
    "context"
    "github.com/go-redis/redis/v8"
    "log"
    "os"
    "strconv"
)

var (
    Rdb *redis.Client
    Ctx = context.Background()
)

func InitRedis() {
    addr := os.Getenv("REDIS_ADDR")
    password := os.Getenv("REDIS_PASSWORD")
    dbStr := os.Getenv("REDIS_DB")

    db, err := strconv.Atoi(dbStr)
    if err != nil {
        log.Fatalf("Invalid REDIS_DB value: %v", err)
    }

    Rdb = redis.NewClient(&redis.Options{
        Addr:     addr,
        Password: password,
        DB:       db,
    })

    if err := Rdb.Ping(Ctx).Err(); err != nil {
        log.Fatalf("Could not connect to Redis: %v", err)
    }
}
