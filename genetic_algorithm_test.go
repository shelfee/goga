package goga_test

import (
	"github.com/tomcraven/goga"
	. "gopkg.in/check.v1"

	// "fmt"
	"math/rand"
	"sync"
	"time"
)

const kNumThreads = 4

type GeneticAlgorithmSuite struct {
}

var _ = Suite(&GeneticAlgorithmSuite{})

func helperGenerateExitFunction(numIterations int) func(goga.Genome) bool {
	totalIterations := 0
	return func(goga.Genome) bool {
		totalIterations++
		if totalIterations >= numIterations {
			return true
		}
		return false
	}
}

func (s *GeneticAlgorithmSuite) TestShouldSimulateUntil(t *C) {

	callCount := 0
	exitFunc := func(g goga.Genome) bool {
		callCount++
		return true
	}

	genAlgo := goga.NewGeneticAlgorithm()
	genAlgo.Init(goga.PopulationSize(1), goga.ParallelSimulations(kNumThreads))
	ret := genAlgo.SimulateUntil(exitFunc)
	t.Assert(ret, IsTrue)
	t.Assert(callCount, Equals, 1)

	callCount = 0
	exitFunc2 := func(g goga.Genome) bool {
		callCount++
		if callCount >= 2 {
			return true
		}
		return false
	}
	ret = genAlgo.SimulateUntil(exitFunc2)
	t.Assert(ret, IsTrue)
	t.Assert(callCount, Equals, 2)
}

func (s *GeneticAlgorithmSuite) TestShouldCallMaterAppropriately_1(t *C) {

	numCalls1 := 0
	mateFunc1 := func(a, b goga.Genome) (goga.Genome, goga.Genome) {
		numCalls1++
		return a, b
	}

	numCalls2 := 0
	mateFunc2 := func(a, b goga.Genome) (goga.Genome, goga.Genome) {
		numCalls2++
		return a, b
	}

	m := goga.NewMater(
		[]goga.MaterFunctionProbability{
			{P: 0.5, F: mateFunc1},
			{P: 0.75, F: mateFunc2},
		},
	)

	genAlgo := goga.NewGeneticAlgorithm()
	genAlgo.Init(goga.PopulationSize(2), goga.ParallelSimulations(kNumThreads))
	genAlgo.Mater = m

	numIterations := 1000
	ret := genAlgo.SimulateUntil(helperGenerateExitFunction(numIterations))
	t.Assert(ret, IsTrue)

	sixtyPercent := (numIterations / 100) * 60
	fourtyPercent := (numIterations / 100) * 40
	t.Assert(numCalls1 < sixtyPercent, IsTrue, Commentf("Num calls [%v] percent [%v]", numCalls1, sixtyPercent))
	t.Assert(numCalls1 > fourtyPercent, IsTrue, Commentf("Num calls [%v] percent [%v]", numCalls1, fourtyPercent))

	sixtyFivePercent := (numIterations / 100) * 65
	eightyFivePercent := (numIterations / 100) * 85
	t.Assert(numCalls2 < eightyFivePercent, IsTrue, Commentf("Num calls [%v] percent [%v]", numCalls2, sixtyPercent))
	t.Assert(numCalls2 > sixtyFivePercent, IsTrue, Commentf("Num calls [%v] percent [%v]", numCalls2, fourtyPercent))
}

func (s *GeneticAlgorithmSuite) TestShouldCallMaterAppropriately_2(t *C) {

	numCalls := 0
	mateFunc := func(a, b goga.Genome) (goga.Genome, goga.Genome) {
		numCalls++
		return a, b
	}

	m := goga.NewMater(
		[]goga.MaterFunctionProbability{
			{P: 1, F: mateFunc},
		},
	)

	genAlgo := goga.NewGeneticAlgorithm()
	populationSize := 100
	genAlgo.Init(goga.PopulationSize(populationSize), goga.ParallelSimulations(kNumThreads))
	genAlgo.Mater = m

	numIterations := 1000
	genAlgo.SimulateUntil(helperGenerateExitFunction(numIterations))

	expectedNumIterations := (numIterations * (populationSize / 2))
	expectedNumIterations -= (populationSize / 2)
	t.Assert(numCalls, Equals, expectedNumIterations)
}

func (s *GeneticAlgorithmSuite) TestShouldCallMaterAppropriately_OddSizedPopulation(t *C) {

	numCalls := 0
	mateFunc := func(a, b goga.Genome) (goga.Genome, goga.Genome) {
		numCalls++
		return a, b
	}

	m := goga.NewMater(
		[]goga.MaterFunctionProbability{
			{P: 1, F: mateFunc},
		},
	)

	genAlgo := goga.NewGeneticAlgorithm()
	populationSize := 99
	genAlgo.Init(goga.PopulationSize(populationSize), goga.ParallelSimulations(kNumThreads))
	genAlgo.Mater = m

	numIterations := 1000
	genAlgo.SimulateUntil(helperGenerateExitFunction(numIterations))

	resultantPopulationSize := len(genAlgo.GetPopulation())
	t.Assert(resultantPopulationSize, Equals, populationSize)
}

type MyEliteConsumerCounter struct {
	NumCalls int
}

func (ec *MyEliteConsumerCounter) OnElite(g goga.Genome) {
	ec.NumCalls++
}

func (s *GeneticAlgorithmSuite) TestShouldCallIntoEliteConsumer(t *C) {

	ec := MyEliteConsumerCounter{}
	genAlgo := goga.NewGeneticAlgorithm()
	genAlgo.Init(goga.PopulationSize(1), goga.ParallelSimulations(kNumThreads))
	genAlgo.EliteConsumer = &ec

	numIterations := 42
	ret := genAlgo.SimulateUntil(helperGenerateExitFunction(numIterations))
	t.Assert(ret, IsTrue)
	t.Assert(ec.NumCalls, Equals, numIterations)
}

func (s *GeneticAlgorithmSuite) TestShouldNotSimulateWithNoPopulation(t *C) {

	genAlgo := goga.NewGeneticAlgorithm()

	callCount := 0
	exitFunc := func(g goga.Genome) bool {
		callCount++
		return true
	}
	ret := genAlgo.SimulateUntil(exitFunc)

	t.Assert(ret, IsFalse)
	t.Assert(callCount, Equals, 0)

	genAlgo.Init(goga.PopulationSize(0), goga.ParallelSimulations(kNumThreads))
	ret = genAlgo.SimulateUntil(exitFunc)
	t.Assert(ret, IsFalse)
	t.Assert(callCount, Equals, 0)

	genAlgo.Init(goga.PopulationSize(1), goga.ParallelSimulations(kNumThreads))
	ret = genAlgo.SimulateUntil(exitFunc)
	t.Assert(ret, IsTrue)
	t.Assert(callCount, Equals, 1)
}

func (s *GeneticAlgorithmSuite) TestShouldGetPopulation(t *C) {

	genAlgo := goga.NewGeneticAlgorithm()

	t.Assert(genAlgo.GetPopulation(), HasLen, 0)

	genAlgo.Init(goga.PopulationSize(1), goga.ParallelSimulations(kNumThreads))
	pop := genAlgo.GetPopulation()
	t.Assert(pop, HasLen, 1)

	g := goga.NewGenome(goga.Bitset{})
	t.Assert(pop[0], FitsTypeOf, g)

	genAlgo.Init(goga.PopulationSize(123), goga.ParallelSimulations(kNumThreads))
	t.Assert(genAlgo.GetPopulation(), HasLen, 123)

	p1 := genAlgo.GetPopulation()
	p2 := genAlgo.GetPopulation()
	t.Assert(len(p1), Equals, len(p2))
	for i := 0; i < len(p1); i++ {
		t.Assert(p1[i], Equals, p2[i])
	}
}

type MySimulatorCounter struct {
	NumCalls int
	m        sync.Mutex
}

func (ms *MySimulatorCounter) Simulate(goga.Genome) {
	ms.m.Lock()
	ms.NumCalls++
	ms.m.Unlock()
}
func (ms *MySimulatorCounter) OnBeginSimulation() []goga.Genome {
	return nil
}
func (ms *MySimulatorCounter) OnEndSimulation([]goga.Genome) {
}
func (ms *MySimulatorCounter) ExitFunc(goga.Genome) bool {
	return false
}

func (s *GeneticAlgorithmSuite) TestShouldSimulatePopulatonCounter(t *C) {
	genAlgo := goga.NewGeneticAlgorithm()

	ms := MySimulatorCounter{}
	genAlgo.Simulator = &ms
	t.Assert(ms.NumCalls, Equals, 0)

	populationSize := 100
	genAlgo.Init(goga.PopulationSize(populationSize), goga.ParallelSimulations(kNumThreads))

	numIterations := 10
	genAlgo.SimulateUntil(helperGenerateExitFunction(numIterations))
	t.Assert(ms.NumCalls, Equals, numIterations*populationSize)
}

type MySimulatorFitness struct {
	NumIterations     int
	LargestFitnessess []int

	currentLargestFitness int
	m                     sync.Mutex
}

func (ms *MySimulatorFitness) Simulate(g goga.Genome) {
	ms.m.Lock()
	randomFitness := rand.Intn(1000)
	if randomFitness > ms.currentLargestFitness {
		ms.currentLargestFitness = randomFitness
	}
	g.SetFitness(float64(randomFitness))
	ms.m.Unlock()
}
func (ms *MySimulatorFitness) OnBeginSimulation() []goga.Genome {
	ms.currentLargestFitness = 0
	return nil
}
func (ms *MySimulatorFitness) OnEndSimulation([]goga.Genome) {
	ms.LargestFitnessess = append(ms.LargestFitnessess, ms.currentLargestFitness)
}
func (ms *MySimulatorFitness) ExitFunc(goga.Genome) bool {
	return false
}

type MyEliteConsumerFitness struct {
	EliteFitnesses []int
}

func (ec *MyEliteConsumerFitness) OnElite(g goga.Genome) {
	ec.EliteFitnesses = append(ec.EliteFitnesses, int(g.GetFitness()))
}

func (s *GeneticAlgorithmSuite) TestShouldSimulatePopulationAndPassEliteToConsumer(t *C) {
	genAlgo := goga.NewGeneticAlgorithm()

	numIterations := 100
	ms := MySimulatorFitness{NumIterations: numIterations}
	genAlgo.Simulator = &ms

	ec := MyEliteConsumerFitness{}
	genAlgo.EliteConsumer = &ec

	genAlgo.Init(goga.PopulationSize(100), goga.ParallelSimulations(kNumThreads))

	genAlgo.SimulateUntil(helperGenerateExitFunction(numIterations))

	t.Assert(ec.EliteFitnesses, DeepEquals, ms.LargestFitnessess)
}

type MySimulatorOrder struct {
	Order []int

	BeginCalled    bool
	SimulateCalled bool
	EndCalled      bool
}

func (ms *MySimulatorOrder) OnBeginSimulation() []goga.Genome {
	ms.Order = append(ms.Order, 1)
	ms.SimulateCalled = true
	return nil
}
func (ms *MySimulatorOrder) Simulate(g goga.Genome) {
	ms.Order = append(ms.Order, 2)
	ms.BeginCalled = true
}
func (ms *MySimulatorOrder) OnEndSimulation([]goga.Genome) {
	ms.Order = append(ms.Order, 3)
	ms.EndCalled = true
}
func (ms *MySimulatorOrder) ExitFunc(goga.Genome) bool {
	return false
}

func (s *GeneticAlgorithmSuite) TestShouldCallOnBeginEndSimulation(t *C) {
	genAlgo := goga.NewGeneticAlgorithm()

	ms := MySimulatorOrder{}
	genAlgo.Simulator = &ms

	t.Assert(ms.BeginCalled, Equals, false)
	t.Assert(ms.SimulateCalled, Equals, false)
	t.Assert(ms.Order, HasLen, 0)

	genAlgo.Init(goga.PopulationSize(1), goga.ParallelSimulations(kNumThreads))
	genAlgo.SimulateUntil(helperGenerateExitFunction(1))

	// Sleep and give time for threads to start up
	time.Sleep(100 * time.Millisecond)

	t.Assert(ms.BeginCalled, Equals, true)
	t.Assert(ms.SimulateCalled, Equals, true)
	t.Assert(ms.Order, HasLen, 3)
	t.Assert(ms.Order, DeepEquals, []int{1, 2, 3})
}

func (s *GeneticAlgorithmSuite) TestShouldPassEliteToExitFunc(t *C) {
	genAlgo := goga.NewGeneticAlgorithm()

	numIterations := 10
	ms := MySimulatorFitness{NumIterations: numIterations}
	genAlgo.Simulator = &ms

	ec := MyEliteConsumerFitness{}
	genAlgo.EliteConsumer = &ec

	populationSize := 10
	genAlgo.Init(goga.PopulationSize(populationSize), goga.ParallelSimulations(kNumThreads))

	passedGenomeFitnesses := make([]int, populationSize)
	callCount := 0
	exitFunc := func(g goga.Genome) bool {
		passedGenomeFitnesses[callCount] = int(g.GetFitness())

		callCount++
		if callCount >= numIterations {
			return true
		}
		return false
	}

	genAlgo.SimulateUntil(exitFunc)

	t.Assert(passedGenomeFitnesses, DeepEquals, ms.LargestFitnessess)
}

func (s *GeneticAlgorithmSuite) TestShouldNotCallMaterWithGenomesFromPopulation(t *C) {

	genAlgo := goga.NewGeneticAlgorithm()

	mateFunc := func(a, b goga.Genome) (goga.Genome, goga.Genome) {
		population := genAlgo.GetPopulation()
		aFound, bFound := false, false
		for i := range population {
			if a == population[i] {
				aFound = true
				if aFound && bFound {
					break
				}
			} else if b == population[i] {
				bFound = true
				if aFound && bFound {
					break
				}
			}
		}
		t.Assert(aFound, IsFalse)
		t.Assert(bFound, IsFalse)
		return a, b
	}

	m := goga.NewMater(
		[]goga.MaterFunctionProbability{
			{P: 1, F: mateFunc},
		},
	)

	genAlgo.Init(goga.PopulationSize(10), goga.ParallelSimulations(kNumThreads))
	genAlgo.Mater = m

	numIterations := 1000
	genAlgo.SimulateUntil(helperGenerateExitFunction(numIterations))
}

type MySelectorCounter struct {
	CallCount int
}

func (ms *MySelectorCounter) Go(genomes []goga.Genome, totalFitness float64) goga.Genome {
	ms.CallCount++
	return genomes[0]
}

func (s *GeneticAlgorithmSuite) TestShouldCallSelectorAppropriately(t *C) {

	genAlgo := goga.NewGeneticAlgorithm()

	selector := MySelectorCounter{}
	genAlgo.Selector = &selector

	populationSize := 100
	genAlgo.Init(goga.PopulationSize(populationSize), goga.ParallelSimulations(kNumThreads))
	t.Assert(selector.CallCount, Equals, 0)

	numIterations := 100
	genAlgo.SimulateUntil(helperGenerateExitFunction(numIterations))
	t.Assert(selector.CallCount, Equals, (populationSize*numIterations)-populationSize)
}

type MySelectorPassCache struct {
	PassedGenomes []goga.Genome
}

func (ms *MySelectorPassCache) Go(genomes []goga.Genome, totalFitness float64) goga.Genome {
	randomGenome := genomes[rand.Intn(len(genomes))]
	ms.PassedGenomes = append(ms.PassedGenomes, randomGenome)
	return randomGenome
}

type MyMaterPassCache struct {
	PassedGenomes []goga.Genome
}

func (ms *MyMaterPassCache) Go(a, b goga.Genome) (goga.Genome, goga.Genome) {
	ms.PassedGenomes = append(ms.PassedGenomes, a)
	ms.PassedGenomes = append(ms.PassedGenomes, b)
	return a, b
}
func (ms *MyMaterPassCache) OnElite(goga.Genome) {
}

func (s *GeneticAlgorithmSuite) TestShouldPassSelectedGenomesToMater(t *C) {

	genAlgo := goga.NewGeneticAlgorithm()

	selector := MySelectorPassCache{}
	genAlgo.Selector = &selector

	mater := MyMaterPassCache{}
	genAlgo.Mater = &mater
	genAlgo.Simulator = &MySimulatorFitness{}

	genAlgo.Init(goga.PopulationSize(100), goga.ParallelSimulations(kNumThreads))
	genAlgo.SimulateUntil(helperGenerateExitFunction(100))

	t.Assert(len(mater.PassedGenomes), Equals, len(selector.PassedGenomes))
	t.Assert(mater.PassedGenomes, DeepEquals, selector.PassedGenomes)
}

type MyBitsetCreateCounter struct {
	NumCalls int
}

func (gc *MyBitsetCreateCounter) Go() goga.Bitset {
	gc.NumCalls++
	return goga.Bitset{}
}

func (s *GeneticAlgorithmSuite) TestShouldCallIntoBitsetCreate(t *C) {

	genAlgo := goga.NewGeneticAlgorithm()

	bitsetCreate := MyBitsetCreateCounter{}
	genAlgo.BitsetCreate = &bitsetCreate

	numGenomes := 100
	genAlgo.Init(goga.PopulationSize(numGenomes), goga.ParallelSimulations(kNumThreads))

	t.Assert(bitsetCreate.NumCalls, Equals, numGenomes)
}

type MyMaterPassCache2 struct {
	PassedGenomes  []goga.Genome
	runningFitness float64
}

func (ms *MyMaterPassCache2) Go(a, b goga.Genome) (goga.Genome, goga.Genome) {

	g1, g2 := goga.NewGenome(goga.Bitset{}), goga.NewGenome(goga.Bitset{})

	ms.PassedGenomes = append(ms.PassedGenomes, g1)
	ms.PassedGenomes = append(ms.PassedGenomes, g2)

	g1.SetFitness(ms.runningFitness)
	ms.runningFitness++
	g2.SetFitness(ms.runningFitness)
	ms.runningFitness++

	return g1, g2
}
func (ms *MyMaterPassCache2) OnElite(goga.Genome) {
}

func (s *GeneticAlgorithmSuite) TestShouldReplaceOldPopulationWithMatedOne(t *C) {

	mater := MyMaterPassCache2{}

	genAlgo := goga.NewGeneticAlgorithm()
	genAlgo.Mater = &mater
	populationSize := 10
	genAlgo.Init(goga.PopulationSize(populationSize), goga.ParallelSimulations(kNumThreads))
	genAlgo.SimulateUntil(helperGenerateExitFunction(2))

	genAlgoPopulation := genAlgo.GetPopulation()
	t.Assert(mater.PassedGenomes, HasLen, populationSize)
	t.Assert(genAlgoPopulation, HasLen, populationSize)

	for i := 0; i < populationSize; i++ {
		t.Assert(mater.PassedGenomes[i].GetFitness(), Equals, genAlgoPopulation[i].GetFitness())
		t.Assert(mater.PassedGenomes[i], Equals, genAlgoPopulation[i])
	}
}

type MySimulatorCallTracker struct {
	NumBeginSimulationsUntilExit int

	NumBeginSimulationCalls int
	NumSimulateCalls        int
	m                       sync.Mutex
}

func (ms *MySimulatorCallTracker) Simulate(goga.Genome) {
	ms.m.Lock()
	ms.NumSimulateCalls++
	ms.m.Unlock()
}
func (ms *MySimulatorCallTracker) OnBeginSimulation() []goga.Genome {
	ms.NumBeginSimulationCalls++
	return nil
}
func (ms *MySimulatorCallTracker) OnEndSimulation([]goga.Genome) {
}
func (ms *MySimulatorCallTracker) ExitFunc(goga.Genome) bool {
	return (ms.NumBeginSimulationCalls >= ms.NumBeginSimulationsUntilExit)
}

func (s *GeneticAlgorithmSuite) TestShouldSimulateUsingSimulatorExitFunction(t *C) {
	genAlgo := goga.NewGeneticAlgorithm()

	ms := MySimulatorCallTracker{}
	ms.NumBeginSimulationsUntilExit = 5
	genAlgo.Simulator = &ms
	t.Assert(ms.NumBeginSimulationCalls, Equals, 0)
	t.Assert(ms.NumSimulateCalls, Equals, 0)

	populationSize := 100
	genAlgo.Init(goga.PopulationSize(populationSize), goga.ParallelSimulations(kNumThreads))
	genAlgo.Simulate()

	t.Assert(ms.NumSimulateCalls, Equals, ms.NumBeginSimulationsUntilExit*populationSize)
	t.Assert(ms.NumBeginSimulationCalls, Equals, ms.NumBeginSimulationsUntilExit)
}
