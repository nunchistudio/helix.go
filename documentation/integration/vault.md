The integration uses [the official Go library](https://pkg.go.dev/github.com/hashicorp/vault/api)
maintained by the HashiCorp team.

Install the Go module with:
```sh
$ go get go.nunchi.studio/helix/integration/vault
```

Simple example on how to import, configure, and use the integration:
```go
import (
  "context"
  "fmt"

  "go.nunchi.studio/helix/integration/vault"
  "go.nunchi.studio/helix/service"
)

func main() {
  cfg := vault.Config{
    Address:   "http://127.0.0.1:8200",
    Namespace: "custom",
    Token:     "my_token",
  }

  client, err := vault.Connect(cfg)
  if err != nil {
    return err
  }

  ctx := context.Background()
  kv := client.KeyValue(ctx, "/mountpath")
  secret, err := kv.Get(ctx, "secretpath")
  if err != nil {
    // ...
  }

  fmt.Println("Secret:", secret.Data)

  if err := service.Start(); err != nil {
    panic(err)
  }

  if err := service.Close(); err != nil {
    panic(err)
  }
}
```
