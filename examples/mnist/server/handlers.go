package server

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/patrikeh/go-deep/examples/mnist/server/mnist"
	"github.com/patrikeh/go-deep/examples/mnist/server/utils"
	"github.com/patrikeh/go-deep/training"
)

func Index(e *echo.Echo) echo.HandlerFunc {
	return echo.StaticFileHandler("dist/index.html", e.Filesystem)
}

func Files(e *echo.Echo) echo.HandlerFunc {
	return echo.StaticDirectoryHandler(
		echo.MustSubFS(e.Filesystem, "dist"),
		false,
	)
}

type predictRequest struct {
	Image    string `json:"image"`
	Expected *int   `json:"expected,omitempty"`
}

type predictResponse struct {
	Prediction  int             `json:"prediction"`
	Expected    *int            `json:"expected,omitempty"`
	Correct     *bool           `json:"correct"`
	Predictions map[int]float64 `json:"predictions"`
}

func Predict(c echo.Context) error {
	if neuralNetwork == nil {
		return c.JSON(500, utils.WrapError("neural network not initialized", nil))
	}

	var req predictRequest
	if err := c.Bind(&req); err != nil {
		return err
	}
	img, err := utils.Base64ToImage(req.Image)
	if err != nil {
		return c.JSON(400, utils.WrapError("invalid image", err))
	}
	_ = utils.SaveImage(img, "dist/image.png")

	tensor, err := utils.ImageToTensor(img)
	if err != nil {
		return c.JSON(400, utils.WrapError("could not convert image to tensor", err))
	}

	fmt.Println(mnist.String(tensor))

	prediction := neuralNetwork.Predict(tensor)
	predictedIndex := mnist.Decode(prediction)

	c.Logger().Printf("prediction: %d", predictedIndex)

	resp := predictResponse{
		Prediction:  predictedIndex,
		Predictions: make(map[int]float64),
	}
	if req.Expected != nil {
		resp.Expected = req.Expected
		correct := *req.Expected == predictedIndex
		resp.Correct = &correct
	}

	for i, p := range prediction {
		resp.Predictions[i] = p
	}

	return c.JSON(200, resp)
}

const (
	correctionSet     = "dist/mnist_correction.csv"
	correctionWeights = "dist/correction.json"
)

func Add(c echo.Context) error {
	var req predictRequest
	if err := c.Bind(&req); err != nil {
		return err
	}
	img, err := utils.Base64ToImage(req.Image)
	if err != nil {
		return c.JSON(400, utils.WrapError("invalid image", err))
	}

	tensor, err := utils.ImageToTensor(img)
	if err != nil {
		return c.JSON(400, utils.WrapError("could not convert image to tensor", err))
	}

	out := training.Example{
		Input:    mnist.ToUint8(tensor),
		Response: mnist.OneHot(10, float64(*req.Expected)),
	}

	if err := mnist.Append(out, correctionSet); err != nil {
		return c.JSON(500, utils.WrapError("could not append to training set", err))
	}

	return c.File(correctionSet)
}

func Train(c echo.Context) error {
	if neuralNetwork == nil {
		return c.JSON(500, utils.WrapError("neural network not initialized", nil))
	}
	correction, err := mnist.Examples(correctionSet)
	if err != nil {
		return c.JSON(500, utils.WrapError("could not load correction set", err))
	}
	trainSet, testSet := correction.Split(0.7)

	config := mnist.TrainingConfig{
		Epochs:      1,
		TrainingSet: trainSet,
		TestSet:     testSet,
		Iterations:  25,
		Trainer:     mnist.Trainer(),
	}

	if err := neuralNetwork.Train(config); err != nil {
		return c.JSON(500, utils.WrapError("could not train neural network", err))
	}

	if err := neuralNetwork.Save(correctionWeights); err != nil {
		return c.JSON(500, utils.WrapError("could not save neural network", err))
	}

	return c.File(correctionWeights)
}
