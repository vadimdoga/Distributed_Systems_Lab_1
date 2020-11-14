# Laboratory work #1 on Distributed Systems

## Requirements
1. Create a `config.ini` file in root
```
[ENVIRONMENT]
HIGH_LIMIT=
LOW_LIMIT=
TIMEOUT=
GATEWAY_ADDR=
IP=
PORT=
MongoDbURI=
BASE_PATH="/api/products"
```
2. go run .

## Task
* Each service has a database. At least 1 with SQL DB and 1 with NOSQL DB (https://www.xplenty.com/blog/the-sql-vs-nosql-difference).
* Tasks are distributed, across multiple requests.
* Status endpoint to show how many tasks are currently processing;
* Limit the number of tasks that can be processed concurrently. Return errors for new tasks if no resources are available;
* Service discovery: upon start, services will register themselves with the gateway;
* add a priority system. Some resources should be saved for high priority tasks;
* add timeouts for tasks. Kill a task once it has been processing for too long;
