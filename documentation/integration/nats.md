The integration uses [the official Go library](https://pkg.go.dev/github.com/nats-io/nats.go)
maintained by the NATS / Synadia team.

Install the Go module with:
```sh
$ go get go.nunchi.studio/helix/integration/nats
```

Simple example on how to import, configure, and use the integration:
```go
import (
  "context"

  natsinte "go.nunchi.studio/helix/integration/nats"
  "go.nunchi.studio/helix/service"

  "github.com/nats-io/nats.go"
)

func main() {
  cfg := nats.Config{
    Addresses: []string{"nats://localhost:4222"},
  }

  js, err := nats.Connect(cfg)
  if err != nil {
    return err
  }

  ctx := context.Background()
  js.Publish(ctx, &nats.Msg{
    Subject: "demo",
    Sub: &nats.Subscription{
    	Queue: "demo-queue",
    },
    Data: []byte(`{ "hello": "world" }`),
  })

  if err := service.Start(); err != nil {
    panic(err)
  }

  if err := service.Close(); err != nil {
    panic(err)
  }
}
```
