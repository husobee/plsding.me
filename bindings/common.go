package bindings

import "github.com/labstack/echo"

type Validatable interface {
	Validate(echo.Context) error
}
