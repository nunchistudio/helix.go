If the context provided contains a span then the newly-created span will be a
child of that span, otherwise it will be a root span.

In the example below, we pass the HTTP request's context. The REST router
integration automatically handles tracing. Therefore, the custom span created
will be a child of the HTTP request's span, with no additional work on your end.
```go
import (
  "net/http"

  "go.nunchi.studio/helix/telemetry/trace"
)

router.GET("/path", func(rw http.ResponseWriter, req *http.Request) {
  
  // ...

  _, span := trace.Start(req.Context(), trace.SpanKindServer, "Custom Span")
  defer span.End()

  if 2+2 == 4 {
    span.RecordError("this is a demo error based on a dummy condition", errors.New("any error"))
  }

  // ...

})
```
