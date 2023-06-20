The integration uses [the official Go library](https://pkg.go.dev/github.com/ClickHouse/clickhouse-go/v2)
maintained by the ClickHouse team.

Install the Go module with:
```sh
$ go get go.nunchi.studio/helix/integration/clickhouse
```

Simple example on how to import, configure, and use the integration:
```go
import (
  "context"

  "go.nunchi.studio/helix/integration/clickhouse"
  "go.nunchi.studio/helix/service"
)

func main() {
  cfg := clickhouse.Config{
    Addresses: []string{"127.0.0.1:8123"},
    Database:  "default",
  }

  db, err := clickhouse.Connect(cfg)
  if err != nil {
    return err
  }

  ctx := context.Background()
  rows, err := db.Query(ctx, "QUERY", args...)
  if err != nil {
    // ...
  }

  defer rows.Close()
  for rows.Next() {
    // ...
  }

  if err := service.Start(); err != nil {
    panic(err)
  }

  if err := service.Close(); err != nil {
    panic(err)
  }
}
```
