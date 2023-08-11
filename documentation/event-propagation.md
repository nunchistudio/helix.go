In the example below, we create an `Event` object using the `event` package.
We then create a new `context.Context` by calling `event.ContextWithEvent`. This
returns a new context including the event created.

helix integrations automatically read/write an `Event` from/into a context when
possible. The integration then passes the `Event` in the appropriate headers.
For example, the NATS JetStream integration achieves this by passing and reading
an `Event` from the messages' headers.

```go
import (
  "go.nunchi.studio/helix/event"
)

router.POST("/anything", func(rw http.ResponseWriter, req *http.Request) {
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

  ctx := event.ContextWithEvent(req.Context(), e)
  msg := &nats.Msg{
    Subject: "demo",
    Sub: &nats.Subscription{
      Queue: "demo-queue",
    },
    Data: []byte(`{ "hello": "world" }`),
  }

  js.Publish(ctx, msg)

  handlerfunc.Accepted(rw, req)
})
```
