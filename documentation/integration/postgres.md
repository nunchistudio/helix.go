The integration uses the [jackc/pgx](https://pkg.go.dev/github.com/jackc/pgx/v5)
Go library.

Install the Go module with:
```sh
$ go get go.nunchi.studio/helix/integration/postgres
```

Simple example on how to import, configure, and use the integration:
```go
import (
  "context"

  "go.nunchi.studio/helix/integration/postgres"
  "go.nunchi.studio/helix/service"
)

func main() {
  cfg := postgres.Config{
    Address:  "127.0.0.1:5432",
    Database: "my_db",
    User:     "username",
    Password: "password",
  }

  db, err := postgres.Connect(cfg)
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
