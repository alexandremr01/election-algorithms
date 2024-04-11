# Election Algorithms

## Development

For development, you can open one of the process in interactive mode through `docker compose run -it p1 bash` and then `go run main.go --algorithm raft`.

## Run in Docker

Run with `docker compose up`. 

## Add a New Algorithm

Add a package inside `algorithms`, such as `bully` and `raft`. You must implement a builder from the format 

```go
func(*config.Config, *state.State, *client.Client) types.Algorithm
```

Then, add this builder to the list at `algorithms/algorithms.go`. It will the be used by the rest of the code.