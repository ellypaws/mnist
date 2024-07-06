package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"io"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/patrikeh/go-deep/training"

	deep "github.com/patrikeh/go-deep"
)

/*
mnist classifier
mnist is a set of hand-written digits 0-9
the dataset in a sane format (as used here) can be found at:
https://pjreddie.com/projects/mnist-in-csv/
*/
func main() {
	rand.Seed(time.Now().UnixNano())

	train, err := load("server/dist/mnist_train.csv")
	if err != nil {
		panic(err)
	}
	test, err := load("server/dist/mnist_test.csv")
	if err != nil {
		panic(err)
	}

	for i := range train {
		for j := range train[i].Input {
			train[i].Input[j] = train[i].Input[j] / 255
		}
	}
	for i := range test {
		for j := range test[i].Input {
			test[i].Input[j] = test[i].Input[j] / 255
		}
	}
	test.Shuffle()
	train.Shuffle()

	neural := deep.NewNeural(&deep.Config{
		Inputs:     len(train[0].Input),
		Layout:     []int{50, 10},
		Activation: deep.ActivationReLU,
		Mode:       deep.ModeMultiClass,
		Weight:     deep.NewNormal(0.6, 0.1), // slight positive bias helps ReLU
		Bias:       true,
	})

	//trainer := training.NewBatchTrainer(training.NewSGD(0.01, 0.5, 1e-6, true), 1, 200, 8)
	trainer := training.NewBatchTrainer(training.NewAdam(0.02, 0.9, 0.999, 1e-8), 1, 200, 16)

	fmt.Printf("training: %d, val: %d, test: %d\n", len(train), len(test), len(test))

	expected := deep.ArgMax(test[0].Response)
	start := time.Now()

	green := lipgloss.NewStyle().Foreground(lipgloss.Color("#00ff00"))
	red := lipgloss.NewStyle().Foreground(lipgloss.Color("#ff0000"))

	fmt.Printf("expected: %v\n", expected)
	const epochs = 15
	const iterations = 1
	for epoch := range epochs {
		var numberPrint strings.Builder
		var column int
		for i, in := range test[0].Input {
			if epoch > 0 {
				break
			}
			if i%28 == 0 {
				column++
				numberPrint.WriteString(fmt.Sprintf("\n%2d: ", column))
			}
			numberPrint.WriteString(lipgloss.NewStyle().Foreground(toColor(in)).Render("██"))
		}
		fmt.Print(numberPrint.String())

		prediction := neural.Predict(test[0].Input)
		predictedIndex := deep.ArgMax(prediction)

		out := fmt.Sprintf("\n\nepoch: %d\noutput: %d (%v)\n%v\n", epoch, predictedIndex, predictedIndex == expected, prediction)
		if predictedIndex == expected {
			out = green.Render(out)
		} else {
			out = red.Render(out)
		}
		fmt.Println(out)

		if epoch < epochs-1 {
			trainStart := time.Now()
			trainer.Train(neural, train, test, iterations)
			fmt.Printf("train time: %s/%s\n", time.Since(trainStart), time.Since(start))
		}
	}
}

func toColor(in float64) lipgloss.Color {
	n := int(in * 255)
	return lipgloss.Color(fmt.Sprintf("#%02x%02x%02x", n, n, n))
}

func load(path string) (training.Examples, error) {
	f, err := os.Open(path)
	defer f.Close()
	if err != nil {
		return nil, err
	}
	r := csv.NewReader(bufio.NewReader(f))

	var examples training.Examples
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		examples = append(examples, toExample(record))
	}

	return examples, nil
}

func toExample(in []string) training.Example {
	res, err := strconv.ParseFloat(in[0], 64)
	if err != nil {
		panic(err)
	}
	resEncoded := onehot(10, res)
	var features []float64
	for i := 1; i < len(in); i++ {
		res, err := strconv.ParseFloat(in[i], 64)
		if err != nil {
			panic(err)
		}
		features = append(features, res)
	}

	return training.Example{
		Response: resEncoded,
		Input:    features,
	}
}

// onehot returns a one-hot encoded vector of the given value
// each ordinal value is for each decimal digit 0 to 9
func onehot(classes int, val float64) []float64 {
	res := make([]float64, classes)
	res[int(val)] = 1
	return res
}
