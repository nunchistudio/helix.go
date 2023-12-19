The integration uses [the OpenFeature Go SDK](https://pkg.go.dev/github.com/open-feature/go-sdk/pkg/openfeature)
as well as [the GO Feature Flag provider](https://pkg.go.dev/github.com/open-feature/go-sdk-contrib/providers/go-feature-flag/pkg).

Install the Go module with:
```sh
$ go get go.nunchi.studio/helix/integration/openfeature
```

Simple example on how to import, configure, and use the integration:
```go
import (
  "context"
  "fmt"

  "go.nunchi.studio/helix/event"
  "go.nunchi.studio/helix/integration/openfeature"
  "go.nunchi.studio/helix/service"
)

func main() {
  cfg := openfeature.Config{
    Paths: []string{
      "./flags/marketing.yaml",
      "./flags/canary.yaml",
    },
  }

  ff, err := openfeature.Connect(cfg)
  if err != nil {
    return err
  }

  var e = event.Event{
    Name:   "post.anything",
    UserID: "7469e788-617a-4b6a-8a26-a61f6acd01d3",
    Subscriptions: []event.Subscription{
      {
        CustomerID:  "2658da04-7c8f-4a7e-9ab0-d5d555b8173e",
        PlanID:      "7781028b-eb48-410d-8cae-c36cffed663d",
        Usage:       "api.requests",
        IncrementBy: 1.0,
      },
    },
  }

  ctx := event.ContextWithEvent(context.Background(), e)
  details, err := app.FeatureFlags.EvaluateString(ctx, "flag-name", "default value", e.UserID)
  if err != nil {
    // ...
  }

  fmt.Println("Variant:", details.Variant)
  fmt.Println("Value:", details.Value)

  if err := service.Start(); err != nil {
    panic(err)
  }

  if err := service.Close(); err != nil {
    panic(err)
  }
}
```
