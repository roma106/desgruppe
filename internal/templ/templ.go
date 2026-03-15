package templ

import (
	"io"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

type Renderer struct{}

func (t *Renderer) Render(w io.Writer, name string, data any, c echo.Context) error {

	cmp, ok := data.(templ.Component)
	if !ok {
		return echo.NewHTTPError(500, "Failed to render templ component")
	}
	return cmp.Render(c.Request().Context(), w)
}
