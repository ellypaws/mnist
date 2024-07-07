package deep

import (
	"math"
	"slices"
)

// Mean of xx
func Mean(xx []float64) float64 {
	var sum float64
	for _, x := range xx {
		sum += x
	}
	return sum / float64(len(xx))
}

// Variance of xx
func Variance(xx []float64) float64 {
	if len(xx) == 1 {
		return 0.0
	}
	m := Mean(xx)

	var variance float64
	for _, x := range xx {
		variance += math.Pow((x - m), 2)
	}

	return variance / float64(len(xx)-1)
}

// StandardDeviation of xx
func StandardDeviation(xx []float64) float64 {
	return math.Sqrt(Variance(xx))
}

// Standardize (z-score) shifts distribution to μ=0 σ=1
func Standardize(xx []float64) {
	m := Mean(xx)
	s := StandardDeviation(xx)

	if s == 0 {
		s = 1
	}

	for i, x := range xx {
		xx[i] = (x - m) / s
	}
}

// Normalize scales to (0,1)
func Normalize[f ~float64](xx []f) {
	min, max := slices.Min(xx), slices.Max(xx)
	for i, x := range xx {
		xx[i] = (x - min) / (max - min)
	}
}

// Min is the smallest element
func Min[f ~float64](xx []f) f {
	min := xx[0]
	for _, x := range xx {
		if x < min {
			min = x
		}
	}
	return min
}

// Max is the largest element
func Max[f ~float64](xx []f) f {
	max := xx[0]
	for _, x := range xx {
		if x > max {
			max = x
		}
	}
	return max
}

// ArgMax is the index of the largest element
func ArgMax[f ~float64](xx []f) int {
	max, idx := xx[0], 0
	for i, x := range xx {
		if x > max {
			max, idx = xx[i], i
		}
	}
	return idx
}

// Sgn is signum
func Sgn[f ~float64](x f) f {
	switch {
	case x < 0:
		return -1.0
	case x > 0:
		return 1.0
	}
	return 0
}

// Sum is sum
func Sum[f ~float64](xx []f) (sum f) {
	for _, x := range xx {
		sum += x
	}
	return
}

// Softmax is the softmax function
func Softmax[f ~float64](xx []f) []f {
	out := make([]f, len(xx))
	var sum f
	max := Max(xx)
	for i, x := range xx {
		out[i] = f(math.Exp(float64(x - max)))
		sum += out[i]
	}
	for i := range out {
		out[i] /= sum
	}
	return out
}

// Round to nearest integer
func Round[f ~float64](x f) f {
	return f(math.Floor(float64(x + .5)))
}

// Dot product
func Dot[f ~float64](xx, yy []f) f {
	var p f
	for i := range xx {
		p += xx[i] * yy[i]
	}
	return p
}
