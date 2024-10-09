Install the Go module with:
```sh
$ go get go.nunchi.studio/helix/integration/temporal
```

Simple example on how to import, configure, and use the integration:
```go
import (
  "go.nunchi.studio/helix/integration/temporal"
  "go.nunchi.studio/helix/service"

  "go.temporal.io/sdk/activity"
  "go.temporal.io/sdk/workflow"
)

func main() {
  cfg := temporal.Config{
    Address:   "localhost:7233",
    Namespace: "default",
    Worker: temporal.ConfigWorker{
      Enabled:   true,
      TaskQueue: "demo",
    },
  }

  _, w, err := temporal.Connect(cfg)
  if err != nil {
    return err
  }

  w.RegisterWorkflow(YourWorkflowDefinition, workflow.RegisterOptions{
    Name: "workflow",
  })

  w.RegisterActivity(YourSimpleActivityDefinition, activity.RegisterOptions{
    Name: "activity",
  })

  if err := service.Start(); err != nil {
    panic(err)
  }

  if err := service.Close(); err != nil {
    panic(err)
  }
}
```
