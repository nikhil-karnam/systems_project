A fault-tolerant task processor on AWS EKS (Kubernetes) with a Redis broker, 4 Go worker pods, a Go API, a retry sweeper, and a React UI hosted on S3.

Clicking send on the web page makes the API push 200 tasks with ticket IDs into Redis. Workers pull tasks, lowercase the message, upload it to S3, and push the result to an output queue. The page polls the API every 50 ms and shows results in one column per worker, so you can see each worker's speed and how tasks are spread out.

Race conditions: each task is claimed by only one worker.
Fault tolerance: if a worker dies mid-task, the sweeper requeues the task.
Idempotency: retried tasks skip the S3 upload if it already happened. Verified with 0 duplicate S3 files after killing a worker mid-task.

GitHub Actions builds Docker images, pushes them to ECR, and rolls them out to the cluster on every push.
