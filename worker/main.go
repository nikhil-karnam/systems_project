
package main

import (
	"context"
	"fmt"
	"time"
	"os"
	"math/rand"
	"strings"
	"encoding/json"
	"bytes"
	"strconv"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/aws"
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

	//s3 setup. credentials and region come from the environment, set in worker.yaml
	bucket := os.Getenv("S3_BUCKET")
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		fmt.Println("aws config error:", err)
		return
	}
	s3c := s3.NewFromConfig(cfg)

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

		//change string
		t.Message = strings.ToLower(t.Message)

		//wrap
		res, _ := json.Marshal(t)

		//external task for idempotency check
		//sends a file to Amazon S3 with ticket number
		//a retried task should NOT repeat this step
		key := "tasks/" + strconv.Itoa(t.TicketID) + "-" + strconv.Itoa(rand.Intn(1000000)) + ".json"
		first, _ := rdb.SAdd(ctx, "written", t.TicketID).Result()
		if first == 1 {
			_, err := s3c.PutObject(ctx, &s3.PutObjectInput{
				Bucket: aws.String(bucket),
				Key:    aws.String(key),
				Body:   bytes.NewReader(res),
			})
			if err != nil {
				fmt.Println("s3 error:", err)
				//un-mark it so a failed write can be retried
				rdb.SRem(ctx, "written", t.TicketID)
			} else {
				fmt.Println("wrote to s3:", key)
			}
		} else {
			fmt.Println("skipped s3, already written:", key)
		}

		//simulate the task actually having process time
		time.Sleep(50 * time.Millisecond)
		
		//worker returns the output to redis --> api --> webpage
		finish.Run(ctx, rdb, []string{"output", "processing", "processTime"}, res, task)		
	}
}
