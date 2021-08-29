# background-processing

## Getting started

This is a background jobs service. It is only in charge to perform jobs in a async fashion.

Install it as it follows:

```bash
chmod +x create_docker-image.sh
./create_docker-image.sh
kubectl apply -f build.yaml
```

This application doesn't have any service in front. The main way to comunicate with it is through database.

## Two tasks available

There are only two task available for now:
1. `Jobs`: It executes a task (time.Sleep) and notifies when starts and finishes as well as the status.
2. `JobCanceller`: It is run when a task is dead (more than 10 seconds without finishing).

For demo purposes, every 3 jobs one is going to be stuck (not finished and therefore cancelled).

The classic behavior to observe is this:

```
| id  | object_id | sleep      | status     | created_at          | queued_at           | started_at          | finished_at         | cancelled_at        |
| 350 |        10 | 1000000000 | DONE!      | 2021-08-29 21:21:36 | 2021-08-29 21:21:38 | 2021-08-29 21:21:38 | 2021-08-29 21:21:39 | NULL                |
| 351 |        10 |       NULL | CANCELLED! | 2021-08-29 21:21:36 | 2021-08-29 21:21:38 | 2021-08-29 21:21:38 | 2021-08-29 21:21:50 | 2021-08-29 21:21:50 |
| 352 |        10 | 1000000000 | DONE!      | 2021-08-29 21:21:36 | 2021-08-29 21:21:38 | 2021-08-29 21:21:38 | 2021-08-29 21:21:39 | NULL                |
```

## Model

This background server is split in two components:
1. `Publishers`: Processes that are constantly (every x seconds) watching the database for items that need an action to perform. e.g. Jobs created but never queued or finished.
2. `Consumer`: Processes that perform the actual task. e.g. Train a model

The idea behind is pretty simple. Each Publisher/Consumer has only one work to do. Split to conquer. The `JobCancelledPublisher` is run every X seconds and will only look for jobs that are dead. It will create one Consumer for each job. That consumer have only one job to do: cancell the job.

The publishers and the consumers are connected through a go channel. The publishers add items to the channel and the consumers are reading that channel. Both the entance and the exit of the channel are processed asynchronously. 

There is an important feature for performance. The consumers are processed in a constant and controlled speed. e.g. 50 tasks per second. One of the main problems of background servers is that are really powerful; they might kill databases if you start shooting thousands queries, updates, deletes, etc. So having a regulator is a performance help.
