package server

import (
	"github.com/labstack/echo/v4"
	"github.com/patrikeh/go-deep/examples/mnist/server/mnist"
)

var neuralNetwork *mnist.Neural

func Run(network *mnist.Neural, middlewares []echo.MiddlewareFunc, options ...func(e *echo.Echo)) error {
	e := echo.New()

	for _, m := range middlewares {
		e.Use(m)
	}

	for _, option := range options {
		option(e)
	}

	registerAs(e.GET, pathHandler{
		"/":  handler{Index(e), nil},
		"/*": handler{Files(e), nil},
	})

	v1 := e.Group("/v1")
	registerAs(v1.POST, postHandlers)
	registerAs(v1.PUT, putHandlers)

	if network != nil {
		neuralNetwork = network
	}

	return e.Start(":1323")
}

var postHandlers = pathHandler{
	"/predict": handler{Predict, nil},
	"/train":   handler{Train, nil},
}

var putHandlers = pathHandler{
	"/train": handler{Add, nil},
}

type route = func(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route

type handler struct {
	handler    func(c echo.Context) error
	middleware []echo.MiddlewareFunc
}

type pathHandler = map[string]handler

func registerAs(route route, pathHandler pathHandler) {
	for path, handler := range pathHandler {
		route(path, handler.handler, handler.middleware...)
	}
}
