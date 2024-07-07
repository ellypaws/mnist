package mnist

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/patrikeh/go-deep/examples/mnist/server/types"
	"github.com/patrikeh/go-deep/examples/mnist/server/utils"
	"io"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/patrikeh/go-deep/training"

	deep "github.com/patrikeh/go-deep"
)

type Neural deep.Neural

/*
mnist classifier
mnist is a set of hand-written digits 0-9
the dataset in a sane format (as used here) can be found at:
https://pjreddie.com/projects/mnist-in-csv/
*/
func init() {
	rand.Seed(time.Now().UnixNano())
}

func New(inputSize int) *Neural {
	return (*Neural)(deep.NewNeural(&deep.Config{
		Inputs:     inputSize,
		Layout:     []int{50, 10},
		Activation: deep.ActivationReLU,
		Mode:       deep.ModeMultiClass,
		Weight:     deep.NewNormal(0.6, 0.1), // slight positive bias helps ReLU
		Bias:       true,
	}))
}

func Load(path string) (*Neural, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer f.Close()
	var dump deep.Dump
	if err := json.NewDecoder(f).Decode(&dump); err != nil {
		return nil, err
	}

	neural := (*Neural)(deep.FromDump(&dump))
	return neural, nil
}

type TrainingConfig struct {
	TrainingSet training.Examples
	TestSet     training.Examples
	Iterations  int
	Trainer     training.Trainer
}

func (n *Neural) Save(path string) error {
	dump := n.network().Dump()
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	if err := json.NewEncoder(f).Encode(dump); err != nil {
		return err
	}
	fmt.Printf("saved to: %s\n", f.Name())

	return nil
}

func (n *Neural) Predict(in []types.Tensor) []float64 {
	return n.network().Predict(types.Coerce[types.Tensor, float64](in))
}

func Decode(prediction []float64) int {
	return deep.ArgMax(prediction)
}

func (n *Neural) Train(config TrainingConfig) error {
	for i := range config.TrainingSet {
		for j := range config.TrainingSet[i].Input {
			config.TrainingSet[i].Input[j] = config.TrainingSet[i].Input[j] / 255
		}
	}
	for i := range config.TestSet {
		for j := range config.TestSet[i].Input {
			config.TestSet[i].Input[j] = config.TestSet[i].Input[j] / 255
		}
	}
	fmt.Printf("training: %d, val: %d, test: %d\n", len(config.TrainingSet), len(config.TestSet), len(config.TestSet))

	config.TestSet.Shuffle()
	config.TrainingSet.Shuffle()

	expected := deep.ArgMax(config.TestSet[0].Response)
	start := time.Now()

	green := lipgloss.NewStyle().Foreground(lipgloss.Color("#00ff00"))
	red := lipgloss.NewStyle().Foreground(lipgloss.Color("#ff0000"))

	fmt.Printf("expected: %v\n", expected)

	fmt.Println(utils.String(utils.DataToTensor(config.TestSet[0].Input)))

	prediction := n.network().Predict(config.TestSet[0].Input)
	predictedIndex := deep.ArgMax(prediction)

	out := fmt.Sprintf("\n\noutput: %d (%v)\n%v\n", predictedIndex, predictedIndex == expected, prediction)
	if predictedIndex == expected {
		out = green.Render(out)
	} else {
		out = red.Render(out)
	}
	fmt.Println(out)

	trainStart := time.Now()
	config.Trainer.Train(n.network(), config.TrainingSet, config.TestSet, config.Iterations)
	fmt.Printf("train time: %s/%s\n", time.Since(trainStart), time.Since(start))

	return nil
}

var trainer *training.BatchTrainer

func Trainer() *training.BatchTrainer {
	if trainer != nil {
		return trainer
	}
	//trainer = training.NewBatchTrainer(training.NewSGD(0.01, 0.5, 1e-6, true), 1, 200, 8)
	trainer = training.NewBatchTrainer(training.NewAdam(0.02, 0.9, 0.999, 1e-8), 1, 200, 16)
	return trainer
}

func Append(example training.Example, path string) error {
	if err := assertFile(path); err != nil {
		return err
	}

	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return err
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	var record []string
	record = append(record, strconv.Itoa(deep.ArgMax(example.Response)))
	for _, in := range example.Input {
		record = append(record, fmt.Sprintf("%d", int(in)))
	}

	if err := w.Write(record); err != nil {
		return err
	}

	return nil
}

func assertFile(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		f, err := os.Create(path)
		if err != nil {
			return err
		}

		err = f.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

func (n *Neural) network() *deep.Neural {
	return (*deep.Neural)(n)
}

func Examples(path string) (training.Examples, error) {
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
	resEncoded := OneHot(10, res)
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

// OneHot returns a one-hot encoded vector of the given value
// each ordinal value is for each decimal digit 0 to 9
func OneHot(classes int, val float64) []float64 {
	res := make([]float64, classes)
	res[int(val)] = 1
	return res
}
