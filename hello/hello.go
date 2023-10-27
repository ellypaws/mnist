package main

import (
	"fmt"
	"gorgonia.org/gorgonia"
	"log"
)

type vars struct {
	graph   *gorgonia.ExprGraph
	x, y, z *gorgonia.Node
}

func main() {
	vars := vars{graph: gorgonia.NewGraph()}

	vars.x = gorgonia.NewScalar(vars.graph, gorgonia.Float64, gorgonia.WithName("x"))
	vars.y = gorgonia.NewScalar(vars.graph, gorgonia.Float64, gorgonia.WithName("y"))

	gorgonia.Let(vars.x, 2.0)
	gorgonia.Let(vars.y, 2.5)

	var err error

	if vars.z, err = gorgonia.Add(vars.x, vars.y); err != nil {
		panic(err)
	}

	machine := gorgonia.NewTapeMachine(vars.graph)
	defer machine.Close()

	if err = machine.RunAll(); err != nil {
		log.Fatal(err)
	}

	fmt.Println(vars.z.Value().Data().(float64))
}
