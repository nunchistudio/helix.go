Create a directory `subscriber`, and initialize Go modules. This directory will
hold the source of the NATS subscription service:
```sh
$ mkdir -p ./services/subscriber && cd ./services/subscriber
$ go mod init
```

Just like for the first service, make sure to add the service to the Go workspace.
At the root path of the project, run:
```sh
$ go work use ./services/*
```

In the `subscriber` directory, create a `service.go`. It contains the source code
related to create, start, and close a helix service:
```go
package main

import (
  "context"
  "errors"

  natsinte "go.nunchi.studio/helix/integration/nats"
  "go.nunchi.studio/helix/service"
  "go.nunchi.studio/helix/telemetry/trace"

  "github.com/nats-io/nats.go/jetstream"
)

/*
App holds the different components needed to run our Go service. In this
case, it only holds a NATS JetStream context.
*/
type App struct {
  JetStream natsinte.JetStream
}

/*
app is the instance of App currently running.
*/
var app *App

/*
NewAndStart creates a new helix service and starts it.
*/
func NewAndStart() error {

  // First, create a new NATS JetStream context. We keep empty config but feel
  // free to dive more later for advanced configuration.
  js, err := natsinte.Connect(natsinte.Config{})
  if err != nil {
    return err
  }

  // Build app with the NATS JetStream context created.
  app = &App{
    JetStream: js,
  }
  
  // Create a new stream in NATS JetStream called "demo-stream", for subject "demo".
  stream, _ := js.CreateStream(context.Background(), jetstream.StreamConfig{
    Name:     "demo-stream",
    Subjects: []string{"demo"},
  })

  // Create a new NATS JetStream consumer called "demo-queue".
  consumer, _ := stream.CreateOrUpdateConsumer(context.Background(), jetstream.ConsumerConfig{
    Name: "demo-queue",
  })

  // Create a new, empty context.
  ctx := context.Background()

  // Start consuming messages from the queue "demo-queue" on subject "demo". We
  // pass the empty context previously created. The context in the callback
  // function is a copy of one the passed, but now contains the Event object at
  // the origin of the trace (if any). You can also create your own span, which
  // will be a child span of the trace found in the context (if any). In our case,
  // the context includes Event created during the HTTP request, as well as the
  // trace handled by the REST router. At any point in time, you can record an
  // error in the span, which will be reported back to the root span.
  consumer.Consume(ctx, func(ctx context.Context, msg jetstream.Msg) {
    _, span := trace.Start(ctx, trace.SpanKindConsumer, "Custom Span")
    defer span.End()

    if 2+2 == 4 {
      span.RecordError("this is a demo error based on a dummy condition", errors.New("any error"))
    }

    msg.Ack()
  })

  // Start the service using the helix's service package. Only one helix service
  // must be running per process. This is a blocking operation.
  err = service.Start()
  if err != nil {
    return err
  }

  return nil
}

/*
Close tries to gracefully close the helix service. This will automatically close
all connections of each integration when applicable. You can add other logic as
well here.
*/
func (app *App) Close() error {
  err := service.Close()
  if err != nil {
    return err
  }

  return nil
}
```

Call the `NewAndStart` and `Close` functions from `main.go`:
```go
package main

/*
main is the entrypoint of our app.
*/
func main() {

  // Create and start the service.
  err := NewAndStart()
  if err != nil {
    panic(err)
  }

  // Try to close the service when done.
  err = app.Close()
  if err != nil {
    panic(err)
  }
}
```

We need to ensure Go dependencies are present. At the root path of the project,
run:
```sh
$ go work sync
```
