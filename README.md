run Kubernetes as EKS on AWS.

tell it to run four workers at once plus redis plus an api. that's its job.

created a web page as a file that exists on Amazon S3. your computer's home program (browser) is edge. when you type in a url or computer address in the search bar of edge, it tells edge to send a message to S3. then S3 sends back html from the web page file and edge runs it and displays a page: an input box to type a message, a button to click print, and an output box. this is an interface with which you can ask edge to send more messages to other computers.

clicking the button is coded to make edge send a message to the api. the api pushes 200 tasks of the message along with ticket IDs to redis.

the worker's job is to constantly ask for tasks from redis. when one is available, take the message, convert everything to lowercase, and do two things: upload the message in a file to Amazon S3, and send the message back to redis to an output queue. the tasks include simulated latency. following the initial button click, edge also makes requests every 50 milliseconds for the api to request information from redis's output queue. one fetch polls a single item from the output queue, and updates the results, which are automatically displayed in the output box. the output is split into one column per worker to show how fast each worker processes a message and how tasks are distributed.

also handles three problems to be solved (distributed systems work):
- race condition: two workers request the same task
- fault tolerance: a worker dies mid task and task must be retried
- idempotency: a retried task has a subtask with an irreversible external output that is already completed and should not be repeated

the web page is written with react / javascript.
the api and workers are written with go.
docker is used to wrap (containerize) the api, workers, (and an added retry checker) so that they can be handled by Kubernetes. redis automatically has a wrapper (docker image).

CI/CD is used to automatically build and deploy all the code except Kubernetes manifests every time you push code with git.
