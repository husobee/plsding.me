package handlers

import "github.com/labstack/echo"

func CreateReminder(c echo.Context) error {
	return nil
}

func GetReminder(c echo.Context) error {
	c.Logger().Info("Reminder id is: ", c.Param("id"))
	return nil
}
