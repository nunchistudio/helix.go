Install the Go module with:
```sh
$ go get go.nunchi.studio/helix/integration/bucket
```

Simple example on how to import, configure, and use the integration:
```go
import (
  "context"

  "go.nunchi.studio/helix/integration/bucket"
  "go.nunchi.studio/helix/service"
)

func main() {
  cfg := bucket.Config{
    Driver:    bucket.DriverAWS,
    Bucket:    "my-bucket",
    Subfolder: "path/to/subfolder/",
  }

  b, err := bucket.Connect(cfg)
  if err != nil {
    return err
  }

  ctx := context.Background()
  blob, err := b.Read(ctx, "blob.json")
  if err != nil {
    // ...
  }

  var anything anyType
  err = json.Unmarshal(blob, &anything)
  if err != nil {
    // ...
  }

  if err := service.Start(); err != nil {
    panic(err)
  }

  if err := service.Close(); err != nil {
    panic(err)
  }
}
```
