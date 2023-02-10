//go:build amd64 && (linux || windows || darwin)

package fj4echo

import (
	"fmt"
	"net/http"

	"github.com/bytedance/sonic"
	"github.com/goccy/go-json"
	"github.com/labstack/echo/v4"
)

func New() echo.JSONSerializer {
	return new(sonicJSONSerializer)
}

func Type() SerializerType {
	return SerializerTypeSonic
}

// sonicJSONSerializer implements JSON encoding using github.com/bytedance/sonic.
type sonicJSONSerializer struct{}

// Serialize converts an interface into a json and writes it to the response.
// You can optionally use the indent parameter to produce pretty JSONs.
func (d sonicJSONSerializer) Serialize(c echo.Context, i interface{}, indent string) error {
	enc := sonic.ConfigDefault.NewEncoder(c.Response())
	if indent != "" {
		enc.SetIndent("", indent)
	}
	return enc.Encode(i)
}

// Deserialize reads a JSON from a request body and converts it into an interface.
func (d sonicJSONSerializer) Deserialize(c echo.Context, i interface{}) error {
	err := sonic.ConfigFastest.NewDecoder(c.Request().Body).Decode(i)
	if ute, ok := err.(*json.UnmarshalTypeError); ok {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Unmarshal type error: expected=%v, got=%v, field=%v, offset=%v", ute.Type, ute.Value, ute.Field, ute.Offset)).SetInternal(err)
	} else if se, ok := err.(*json.SyntaxError); ok {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Syntax error: offset=%v, error=%v", se.Offset, se.Error())).SetInternal(err)
	}
	return err
}
