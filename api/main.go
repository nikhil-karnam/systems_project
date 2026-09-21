package main

import (
	"context"
	"fmt"
	"os"
	"net/http"
	"encoding/json"
	"github.com/redis/go-redis/v9"
)

type Task struct {
	TicketID	int    `json:"ticketID"`
	Message 	string `json:"message"`
}

func main() {
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}
	rdb := redis.NewClient(&redis.Options{Addr: redisAddr})
	ctx := context.Background()

	//listen from the web page
	http.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			return
		}
		var body struct {
			Message string `json:"message"`
		}
		json.NewDecoder(r.Body).Decode(&body)
		message := body.Message
		
		for i := 1; i <= 200; i++ {
    		data, _ := json.Marshal(Task{TicketID: i, Message: message,})
			rdb.LPush(ctx, "tasks", data)
		}
	})

	//listen to redis
	http.HandleFunc("/poll", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")

		task, err := rdb.RPop(ctx, "output").Result()
		if err != nil {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		fmt.Fprint(w, task)
	})

	http.ListenAndServe(":8080", nil)
}