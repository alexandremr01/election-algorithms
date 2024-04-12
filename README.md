# Election Algorithms

## Development

For development, you can open one of the process in interactive mode through `docker compose run -it p1 bash` and then `go run main.go --algorithm raft`. After any changes, we do 

```
go fmt ./...
goimports -w -l .
golangci-lint run
```

## Run

It can wither be run as

- One node per docker container, using `docker compose up`;
- One node per process, using `sh run.sh`.

## Add a New Algorithm

Add a package inside `algorithms`, such as `bully` and `raft`. You must implement a builder from the format 

```go
func(*types.Config, *state.State, *client.Client) types.Algorithm
```

Then, add this builder to the list at `algorithms/algorithms.go`. It will the be used by the rest of the code.