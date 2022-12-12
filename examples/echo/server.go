package main

import (
	"github.com/thoas/stats"
	"net/http"
)

// Stats provides response time, status code count, etc.
var Stats = stats.New()

func main() {
	r := echo.New()

	r.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			beginning, recorder := Stats.Begin(ctx.Response().Writer)
			err := next(ctx)
			Stats.End(beginning, stats.WithRecorder(recorder))
			return err
		}
	})

	r.GET("/stats", func(ctx echo.Context) error {
		return ctx.JSON(http.StatusOK, Stats.Data())
	})

	r.GET("/", func(ctx echo.Context) error {
		return ctx.JSON(http.StatusOK, map[string]string{"hello": "world"})
	})

	r.Start("0.0.0.0:8080")
}
