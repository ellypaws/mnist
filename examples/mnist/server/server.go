package server

import (
	"github.com/labstack/echo/v4"
	"github.com/patrikeh/go-deep/examples/mnist/server/mnist"
)

var neuralNetwork *mnist.Neural

func Run(network *mnist.Neural, middlewares ...echo.MiddlewareFunc) error {
	e := echo.New()

	for _, m := range middlewares {
		e.Use(m)
	}

	registerAs(e.GET, pathHandler{
		"/":  handler{Index(e), nil},
		"/*": handler{Files(e), nil},
	})
	registerAs(e.POST, postHandlers)
	registerAs(e.PUT, putHandlers)

	if network != nil {
		neuralNetwork = network
	}

	return e.Start(":1323")
}

var postHandlers = pathHandler{
	"/v1/predict": handler{Predict, nil},
	"/v1/train":   handler{Train, nil},
}

var putHandlers = pathHandler{
	"/v1/train": handler{Add, nil},
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
