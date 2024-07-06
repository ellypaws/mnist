package server

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/patrikeh/go-deep/examples/mnist/server/mnist"
	"github.com/patrikeh/go-deep/training"
	"image"
	"image/color"
	"image/png"
	"os"
	"strings"
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
		return c.JSON(500, wrapError("neural network not initialized", nil))
	}

	var req predictRequest
	if err := c.Bind(&req); err != nil {
		return err
	}
	img, err := base64ToImage(req.Image)
	if err != nil {
		return c.JSON(400, wrapError("invalid image", err))
	}
	_ = saveImage(img, "dist/image.png")

	tensor, err := imageToTensor(img)
	if err != nil {
		return c.JSON(400, wrapError("could not convert image to tensor", err))
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
	img, err := base64ToImage(req.Image)
	if err != nil {
		return c.JSON(400, wrapError("invalid image", err))
	}

	tensor, err := imageToTensor(img)
	if err != nil {
		return c.JSON(400, wrapError("could not convert image to tensor", err))
	}

	out := training.Example{
		Input:    mnist.ToUint8(tensor),
		Response: mnist.OneHot(10, float64(*req.Expected)),
	}

	if err := mnist.Append(out, correctionSet); err != nil {
		return c.JSON(500, wrapError("could not append to training set", err))
	}

	return c.File(correctionSet)
}

func Train(c echo.Context) error {
	if neuralNetwork == nil {
		return c.JSON(500, wrapError("neural network not initialized", nil))
	}

	config := mnist.TrainingConfig{
		Epochs:      10,
		TrainingSet: correctionSet,
		TestSet:     correctionSet,
		Iterations:  1,
	}

	if err := neuralNetwork.Train(config); err != nil {
		return c.JSON(500, wrapError("could not train neural network", err))
	}

	if err := neuralNetwork.Save(correctionWeights); err != nil {
		return c.JSON(500, wrapError("could not save neural network", err))
	}

	return c.File(correctionWeights)
}

func wrapError(message string, err error) map[string]string {
	if err == nil {
		return nil
	}
	return map[string]string{"error": message, "debug": err.Error()}
}

func base64ToImage(base64Str string) (image.Image, error) {
	base64Str = strings.TrimPrefix(base64Str, "data:image/png;base64,")
	data, err := base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		return nil, err
	}

	img, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	return img, nil
}

func saveImage(img image.Image, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return png.Encode(f, img)
}

func imageToTensor(img image.Image) ([]float64, error) {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	if width != 28 || height != 28 {
		return nil, errors.New("image must be 28x28 pixels")
	}

	var tensor []float64
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			c := img.At(x, y)
			grayColor := color.GrayModel.Convert(c).(color.Gray)
			tensor = append(tensor, normalize(grayColor.Y))
		}
	}

	return tensor, nil
}

func normalize(color uint8) float64 {
	return float64(color) / 255.0
}
