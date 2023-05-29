package main

import (
	"fmt"
	"github.com/tomcraven/goga"
	"math"
	"math/rand"
	"time"
)

type funcMaterSimulator struct {
	function      func([]float64) float64
	paramSize     int
	findSmallest  bool
	transEndValue float64
	transMaxValue float64
	powValue      float64
}

func (sms *funcMaterSimulator) TransToPos(v float64) float64 {
	if v < 0 {
		return 0
	} else {
		return v
	}
	//if v < sms.transEndValue {
	//	return 0
	//}
	//t := sms.transMaxValue / (1 + math.Pow(sms.powValue, -v))
	//return t
}

func (sms *funcMaterSimulator) OnBeginSimulation() {
}
func (sms *funcMaterSimulator) OnEndSimulation() {
}

func (sms *funcMaterSimulator) Simulate(g goga.Genome) {
	params := goga.ParseBitsToFloat64Arr(g.GetBits())
	flag := 1.
	if sms.findSmallest {
		flag = -1.
	}
	g.SetOrigin(sms.function(params))
	g.SetFitness(sms.TransToPos(flag * g.GetOrigin()))
}
func (sms *funcMaterSimulator) ExitFunc(g goga.Genome) bool {
	return false
}

type myBitsetCreate struct {
	paramsSize int
	initValue  []float64
}

func (bc *myBitsetCreate) Go() goga.Bitset {
	b := goga.Bitset{}
	b.Create(bc.paramsSize * 8)
	for i := 0; i < bc.paramsSize; i++ {
		var param float64
		if i < len(bc.initValue) {
			param = bc.initValue[i]
		} else {
			param = rand.ExpFloat64()
			if rand.Intn(2) == 0 {
				param *= -1
			}
		}
		byteArr := goga.Float64ToByte(param)
		for idx := 0; idx < 8; idx++ {
			b.Set(i*8+idx, int(byteArr[idx]))
		}
	}
	return b
}

type myEliteConsumer struct {
	currentIter int
	paramSize   int
}

func (ec *myEliteConsumer) OnElite(g goga.Genome) {
	gBits := g.GetBits()
	ec.currentIter++
	params := goga.ParseBitsToFloat64Arr(gBits)
	fmt.Println(ec.currentIter, "\t", params, "\tfunc value: ", g.GetOrigin(), "\tfitness: ", g.GetFitness())
}

func main() {
	paramSize := 4
	s := funcMaterSimulator{
		paramSize: paramSize,
		function: func(float64s []float64) float64 {
			sum := 0.
			for i := 0; i < len(float64s); i++ {
				sum += -float64(i)*math.Pow(float64s[i], 4) + float64(i)*math.Pow(float64s[i], 3) + math.Pow(float64s[i], 2) + math.Pow(float64s[i], 1)
			}
			return sum
		},
		findSmallest:  false,
		transEndValue: -500,
		transMaxValue: 100000000,
		powValue:      1.001,
	}

	genAlgo := goga.NewGeneticAlgorithm()

	genAlgo.Simulator = &s
	genAlgo.BitsetCreate = &myBitsetCreate{paramsSize: 4}
	genAlgo.EliteConsumer = &myEliteConsumer{}
	mater := goga.FloatMater{
		Precision: 2,
		MaxValue:  1000,
		MinValue:  -1000,
	}
	genAlgo.Mater = goga.NewMater(
		[]goga.MaterFunctionProbability{
			{P: 0.1, F: mater.ArithmeticMutate},
			{P: 1.0, F: mater.ArithmeticExchange},
			{P: 1.0, F: mater.ArithmeticCrossover, UseElite: true},
		},
	)
	genAlgo.Selector = goga.NewSelector(
		[]goga.SelectorFunctionProbability{
			{P: 1.0, F: goga.Roulette},
		},
	)
	populationSize := 500
	numThreads := 10

	genAlgo.Init(goga.PopulationSize(populationSize), goga.ParallelSimulations(numThreads), goga.MaterExtraRatio(4))

	startTime := time.Now()
	genAlgo.Simulate()
	fmt.Println(time.Since(startTime))

	//bitset := b.Go()
	//params := goga.ParseBitsToFloat64Arr(&bitset)
	//for i, p := range params {
	//	fmt.Printf("idx: %d, param: %f\n", i, p)
	//}

}
