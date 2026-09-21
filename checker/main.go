package main

import (
	"context"
	"fmt"
	"strconv"
	"time"
	"os"
	"github.com/redis/go-redis/v9"
)

const tolerance = 5

func main() {
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}
	rdb := redis.NewClient(&redis.Options{Addr: redisAddr})
	ctx := context.Background()

	for {
		time.Sleep(5 * time.Second)

		claimed, err := rdb.HGetAll(ctx, "processTime").Result()
		if err != nil {
			continue
		}

		for task, tsStr := range claimed {
			ts, _ := strconv.ParseInt(tsStr, 10, 64)
			if time.Now().Unix() - ts < tolerance {
				continue
			}

			rdb.LRem(ctx, "processing", 1, task)
			rdb.HDel(ctx, "processTime", task)
			rdb.LPush(ctx, "tasks", task)

			fmt.Println("requeued:", task)
		}
	}
}