package model

import "github.com/labstack/echo/v4"

type EchoGroup struct {
	API     *echo.Group
	Private *echo.Group
	Public  *echo.Group
}
