By passing a Go context, the logger is aware if the log is part of a trace/span.
If so, `trace_id` and `span_id` are added to the log so it can be linked to the
respective trace/span.

In the example below, we pass the HTTP request's context. The REST router
integration automatically handles tracing. Therefore, the log will be associated
to the trace/span with no additional work on your end.
```go
import (
  "net/http"

  "go.nunchi.studio/helix/telemetry/log"
)

router.GET("/path", func(rw http.ResponseWriter, req *http.Request) {
  
  // ...

  log.Warn(req.Context(), "this is a warning")

  // ...

})
```
