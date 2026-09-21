
package main

import (
	"context"
	"fmt"
	"time"
	"os"
	"encoding/json"
	"github.com/redis/go-redis/v9"
)

type Task struct {
	TicketID   int    `json:"ticketID"`
	Message		string `json:"message"`
	WorkerID 	string `json:"workerID"`
}

/*these commands are sent to redis as a package and run by redis itself. thus the worker dying
during the script does not affect the completion of the script. if the worker were to execute
the commands instead and it dies during, the commands remain half executed.*/
var claim = redis.NewScript(`
	local task = redis.call('RPOPLPUSH', KEYS[1], KEYS[2])
	if task then
		redis.call('HSET', KEYS[3], task, ARGV[1])
	end
	return task
`)

var finish = redis.NewScript(`
	redis.call('LPUSH', KEYS[1], ARGV[1])
	redis.call('LREM', KEYS[2], 1, ARGV[2])
	redis.call('HDEL', KEYS[3], ARGV[2])
	return 1
`)

func main() {
	workerID, _ := os.Hostname()
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}
	rdb := redis.NewClient(&redis.Options{Addr: redisAddr})
	ctx := context.Background()

	for {
		task, err := claim.Run(ctx, rdb, []string{"tasks", "processing", "processTime"}, time.Now().Unix()).Text()
		if err != nil {
			time.Sleep(50 * time.Millisecond)
			continue
		}

		//unwrap
		var t Task
		json.Unmarshal([]byte(task), &t)
		t.WorkerID = workerID

		//do shit
		fmt.Printf("%d: %s\n", t.TicketID, t.Message)

		//simulate the task actually having latency
		time.Sleep(50 * time.Millisecond)

		//wrap
		res, _ := json.Marshal(t)
		finish.Run(ctx, rdb, []string{"output", "processing", "processTime"}, res, task)		
	}
}