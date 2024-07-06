package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/patrikeh/go-deep/examples/mnist/server"
	"github.com/patrikeh/go-deep/examples/mnist/server/mnist"
	"os"
	"time"
)

const (
	inputSize = 28 * 28
	weights   = "dist/weights.json"

	epochs      = 1
	iterations  = 25
	trainingSet = "dist/mnist_train.csv"
	testSet     = "dist/mnist_test.csv"

	correctionSet = "dist/mnist_correction.csv"
)

func main() {
	network, err := mnist.Load(weights)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("initializing network")
			network = initialize()
		} else {
			panic(err)
		}
	} else {
		fmt.Println("loaded network")
	}

	if err := server.Run(network, middlewares...); err != nil {
		panic(err)
	}
}

func initialize() *mnist.Neural {
	network := mnist.New(inputSize)
	err := network.Train(mnist.TrainingConfig{
		Epochs:      epochs,
		TrainingSet: trainingSet,
		TestSet:     testSet,
		Iterations:  iterations,
		Trainer:     mnist.Trainer(),
	})
	if err != nil {
		panic(err)
	}

	if err := network.Save(weights); err != nil {
		panic(err)
	}

	return network
}

var middlewares = []echo.MiddlewareFunc{
	middleware.LoggerWithConfig(
		middleware.LoggerConfig{
			Skipper:          nil,
			Format:           `${time_custom}     	${status} ${method}  ${host}${uri} in ${latency_human} from ${remote_ip} ${error}` + "\n",
			CustomTimeFormat: time.DateTime,
		},
	),
	middleware.RemoveTrailingSlash(),
	middleware.Gzip(),
	middleware.Decompress(),
	middleware.NonWWWRedirect(),
	middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}),
}
