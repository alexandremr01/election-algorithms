# Election Algorithms

Extensible comparison of election algorithms. Currently we have implemented raft and bully.

## Running

It can wither be run through Docker or locally.

- One node per docker container, using `docker compose up`;

Locally the step by step is:

1. Initialize one node per process, using `bash run_process.sh bully`.
2. Kill processes with `bash kill_process.sh 4`.
3. Start processes again with `bash run_process.sh bully 4`.
4. Kill remaining processes with `bash kill_process.sh all`.

## Development

For development, you can open one of the processes in interactive mode through `docker compose run -it p1 bash` and then `go run main.go --algorithm raft`. After any changes, we do 

```
go fmt ./...
goimports -w -l .
golangci-lint run
```

## Adding a New Algorithm

Add a package inside `algorithms`, such as `bully` and `raft`. You must implement a builder in the format 

```go
func(*types.Config, *state.State, *client.Client) types.Algorithm
```

Then, add this builder to the list at `algorithms/algorithms.go`. It will the be used by the rest of the code.