package main

import (
	"github.com/tomcraven/goga"
	fo "github.com/tomcraven/goga/function_optimizer"
	"math"
)

func main() {
	paramSize := 4
	requirement := goga.Float64Requirement{
		Precision: 0.001,
		MaxValue:  10000,
		MinValue:  -10000,
		Specific: map[int]struct {
			Precision float64
			MaxValue  float64
			MinValue  float64
		}{
			1: {
				Precision: 1,
				MinValue:  20,
				MaxValue:  100,
			},
			2: {
				Precision: 0.3,
				MinValue:  20,
				MaxValue:  100,
			},
		},
	}
	function := func(float64s []float64) float64 {
		sum := 0.
		for i := 0; i < len(float64s); i++ {
			sum += -5*float64(i+1)*math.Pow(float64s[i], 4) + float64(i+1)*math.Pow(float64s[i], 3) + 4*float64(i+1)*math.Pow(float64s[i], 2) + math.Pow(float64s[i], 1)
		}
		return sum
	}
	transFunc := func(v float64) float64 {
		if v < -10000000 {
			return 0
		} else {
			return v + 10000001
		}
	}
	algo := fo.NewFuncAlgo(fo.Function(function), fo.ParamSize(paramSize), fo.Requirement(&requirement), fo.TransFunc(transFunc))
	algo.Simulate()
}
