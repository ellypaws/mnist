package server

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/patrikeh/go-deep/server/mnist"
	"github.com/patrikeh/go-deep/server/types"
	"github.com/patrikeh/go-deep/server/utils"
	"github.com/patrikeh/go-deep/training"
	"slices"
)

const (
	correctionSet     = "dist/drawing_data.csv"
	correctionWeights = "dist/correction_weights.json"

	iterations = 100
)

var synthesizer = utils.SyntheticConfig{
	Rotate:    &utils.Values{Min: -5, Max: 5},
	Translate: &utils.Values{Min: -5, Max: 5},
	Zoom:      &utils.Values{Min: 1.05, Max: 1.15},
}

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

	tensor := utils.ImageToTensor(img)

	fmt.Println(utils.String(tensor))

	prediction := neuralNetwork.Predict(tensor)
	predictedIndex := mnist.Decode(prediction)

	if req.Expected != nil {
		if *req.Expected == predictedIndex {
			c.Logger().Infof("prediction: %d", predictedIndex)
		} else {
			c.Logger().Warnf("prediction: %d, expected: %d", predictedIndex, *req.Expected)
		}
	} else {
		c.Logger().Debugf("prediction: %d", predictedIndex)
	}

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

type addRequest struct {
	Image    string `json:"image"`
	Expected *int   `json:"expected,omitempty"`
	Correct  *bool  `json:"correct,omitempty"`
}

func Add(c echo.Context) error {
	var req addRequest
	if err := c.Bind(&req); err != nil {
		return err
	}
	img, err := utils.Base64ToImage(req.Image)
	if err != nil {
		return c.JSON(400, utils.WrapError("invalid image", err))
	}

	bin := utils.ImageToBytes(img)

	out := training.Example{
		Input:    types.Coerce[types.Byte, float64](bin),
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

	syntheticData := synthesizer.Synthesize(correction)

	allData := append(slices.Clone(correction), syntheticData.Rotated...)
	allData = append(allData, syntheticData.Translated...)
	allData = append(allData, syntheticData.Zoomed...)

	trainSet, testSet := allData.Split(0.8)

	config := mnist.TrainingConfig{
		TrainingSet: trainSet,
		TestSet:     testSet,
		Iterations:  iterations,
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
