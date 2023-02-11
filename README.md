# Fast JSON Serialization for Echo

By default, [Echo](https://github.com/labstack/echo) uses Go's stdlib `encoding/json`, which have more performant alternatives such as [jsoniter](https://github.com/json-iterator/go), [go-json](https://github.com/goccy/go-json), and [sonic](https://github.com/bytedance/sonic). In almost all of the cases `encoding/json` is sufficient, and I never needed a faster JSON encoding for my backend. But while playing with Echo, I discovered that it has support for custom JSON serialization (`JSONSerializer` interface), and the idea of combining sonic and go-json seemed gorgeous to me, so I created this package.

fj4echo conditionally selects the JSON library. If processor is amd64, it uses sonic, otherwise it uses go-json.

## Usage

`go get -u github.com/tomruk/fj4echo`

```go
import (
    "fmt"
    "github.com/tomruk/fj4echo"
    "github.com/labstack/echo/v4"
)

func main() {
    e := echo.New()
    e.JSONSerializer = fj4echo.New()

    serializerType := fj4echo.Type()
    fmt.Printf("Serializer: %s\n", serializerType.Name())

    e.Start(":8000")
}
```

You can programmatically check which serializer is being used with `SerializerType` enum (see [serializer_type.go](serializer_type.go)):

```go
serializerType := fj4echo.Type()
switch serializerType {
    case fj4echo.SerializerTypeSonic:
        fmt.Println("sonic is choosen. This means that the processor is amd64")
    case fj4echo.SerializerTypeGoJSON:
        fmt.Println("go-json is choosen")
}
```

## Customization

Serializers can be customized. Keep in mind that you might jeopardize security or cause your backend to consume more memory. Please read and gain an understanding of the settings before you change them:
- [sonic](https://github.com/bytedance/sonic/blob/main/api.go)
- [go-json](https://github.com/goccy/go-json/blob/master/option.go)

fj4echo uses a default configuration, and you can find it and explanations of the choices behind its settings inside the `DefaultConfig` function (in [config.go](config.go)).

The default configuration was written with the following things in mind:
- fj4echo was designed for backend, not for other cases. This means:
    - Do not compromise security.
    - Try not to compromise network latency.

To use a custom configuration, use `NewWithConfig` function:

```go
import (
    "fmt"
    "github.com/tomruk/fj4echo"
    "github.com/labstack/echo/v4"
)

func main() {
    e := echo.New()
    config := fj4echo.Config{
        SonicConfig: sonic.Config{
            // I want to live in a dangerous world, so I enable this.
            EscapeHTML: true,
            // Redundant, and the last thing needed on this Earth.
            SortMapKeys: true,
        },
        GoJSON: GoJSONConfig{
            EncodeOptions: []json.EncodeOptionFunc{
                // I want to live in a dangerous world, so I disable HTML escape.
                json.DisableHTMLEscape(),
            },
        },
    }
    e.JSONSerializer = fj4echo.NewWithConfig(config)

    e.Start(":8000")
}
```