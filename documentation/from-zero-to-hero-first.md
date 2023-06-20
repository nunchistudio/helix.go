We will create a directory called `helix-platform` in your `$GOPATH`, such as:
```sh
/Users/username/go/src/github.com/username/helix-platform
```

Create the directory and initialize a Go workspace:
```sh
$ mkdir ./helix-platform && cd ./helix-platform
$ go work init
```

Create a directory `httpapi`, and initialize Go modules. This directory will hold
the source of the HTTP API service:
```sh
$ mkdir -p ./services/httpapi && cd ./services/httpapi
$ go mod init
```

Make sure to add the service to the Go workspace. At the root path of the project,
run:
```sh
$ go work use ./services/*
```

In the `httpapi` directory, create a `service.go`. It contains the source code
related to create, start, and close a helix service:
```go
package main

import (
  "net/http"

  "go.nunchi.studio/helix/integration/rest"
  "go.nunchi.studio/helix/integration/rest/handlerfunc"
  "go.nunchi.studio/helix/service"
)

/*
App holds the different components needed to run our Go service. In this
case, it only holds a REST router for now.
*/
type App struct {
  REST rest.REST
}

/*
app is the instance of App currently running.
*/
var app *App

/*
NewAndStart creates a new helix service and starts it.
*/
func NewAndStart() error {

  // First, create a new REST router. We keep empty config but feel free to
  // dive more later for configuring OpenAPI behavior.
  router, err := rest.New(rest.Config{})
  if err != nil {
    return err
  }

  // Build app with the router created.
  app = &App{
    REST: router,
  }

  // Add a simple route, returning a 202 HTTP response.
  router.POST("/anything", func(rw http.ResponseWriter, req *http.Request) {
    handlerfunc.Accepted(rw, req)
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

Run the service with:
```sh
$ go run ./services/httpapi
```
