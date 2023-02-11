//go:build !amd64 || (amd4 && !(linux || windows || darwin))

package fj4echo

import (
	"fmt"
	"net/http"

	"github.com/goccy/go-json"
	"github.com/labstack/echo/v4"
)

func New() echo.JSONSerializer {
	defaultConfig := DefaultConfig()
	return &goJSONSerializer{config: defaultConfig.GoJSON}
}

func NewWithConfig(config Config) echo.JSONSerializer {
	return &goJSONSerializer{config: config.GoJSON}
}

func Type() SerializerType {
	return SerializerTypeGoJSON
}

// goJSONSerializer implements JSON encoding using github.com/goccy/go-json
type goJSONSerializer struct {
	config GoJSONConfig
}

// Serialize converts an interface into a json and writes it to the response.
// You can optionally use the indent parameter to produce pretty JSONs.
func (s *goJSONSerializer) Serialize(c echo.Context, i interface{}, indent string) error {
	enc := json.NewEncoder(c.Response())
	if indent != "" {
		enc.SetIndent("", indent)
	}
	return enc.EncodeWithOption(i, s.config.EncodeOptions...)
}

// Deserialize reads a JSON from a request body and converts it into an interface.
func (s *goJSONSerializer) Deserialize(c echo.Context, i interface{}) error {
	err := json.NewDecoder(c.Request().Body).DecodeWithOption(i, s.config.DecodeOptions...)
	if ute, ok := err.(*json.UnmarshalTypeError); ok {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Unmarshal type error: expected=%v, got=%v, field=%v, offset=%v", ute.Type, ute.Value, ute.Field, ute.Offset)).SetInternal(err)
	} else if se, ok := err.(*json.SyntaxError); ok {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Syntax error: offset=%v, error=%v", se.Offset, se.Error())).SetInternal(err)
	}
	return err
}
