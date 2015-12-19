package main

import (
	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
	"jdict/lib/jmdict"
	"net/http"
)

// Handlers
func query(c *echo.Context) error {
	key := c.Param("key")
	result := jmdict.Query(key)
	return c.JSON(http.StatusOK, result)
}

func main() {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(mw.Logger())
	e.Use(mw.Recover())

	// Routes
	e.Get("/query/:key", query)

	// Start server
	e.Run(":3000")
}
