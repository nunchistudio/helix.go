Install the Go module with:
```sh
$ go get go.nunchi.studio/helix/integration/rest
```

Simple example on how to import, configure, and use the integration:
```go
import (
  "net/http"

  "go.nunchi.studio/helix/integration/rest"
  "go.nunchi.studio/helix/service"
)

func main() {
  cfg := rest.Config{
    Address: ":8080",
    OpenAPI: rest.ConfigOpenAPI{
      Enabled:     true,
      Description: "./descriptions/openapi.yaml",
    },
  }

  router, err := rest.New(cfg)
  if err != nil {
    return err
  }

  router.POST("/users/:id", func(rw http.ResponseWriter, req *http.Request) {
    params, ok := rest.ParamsFromContext(req.Context())
    if !ok {
      rest.WriteEmptyNotFound(rw, req)
      return
    }

    userID := params["id"]
    
    // ...
    
    rest.WriteEmptyAccepted(rw, req)
  })

  if err := service.Start(); err != nil {
    panic(err)
  }

  if err := service.Close(); err != nil {
    panic(err)
  }
}
```
