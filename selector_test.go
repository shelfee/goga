package goga_test

import (
	"math"

	"github.com/tomcraven/goga"
	. "gopkg.in/check.v1"
	// "fmt"
)

type SelectorSuite struct {
	selector goga.Selector
}

var _ = Suite(&SelectorSuite{})

func (s *SelectorSuite) TestShouldInstantiate(t *C) {
	// Tested as part of fixture setup
}

func (s *SelectorSuite) TestShouldRoulette(t *C) {

	// for i := 0; i < 100; i++ {
	// 	numGenomes := 100
	// 	genomeArray := make( []goga.IGenome, numGenomes )
	// 	totalFitness := 0
	// 	for i := 0; i < numGenomes; i++ {
	// 		genomeArray[i] = goga.NewGenome( goga.Bitset{} )
	// 		genomeArray[i].SetFitness( i )
	// 		totalFitness += i
	// 	}

	// 	numIterations := 100000
	// 	pickedGenomeFrequency := make( []int, numGenomes )

	// 	for i := 0; i < numIterations; i++ {
	// 		g := s.selector.Roulette( genomeArray, totalFitness )
	// 		pickedGenomeFrequency[g.GetFitness()]++
	// 	}

	// 	comparisonJump := 20
	// 	for i := 0; i < numGenomes - comparisonJump; i++ {
	// 		t.Assert( pickedGenomeFrequency[i] <= pickedGenomeFrequency[i + comparisonJump], IsTrue,
	// 			Commentf( "Piced freq [%v], comparing [%v] < [%v]",
	// 				pickedGenomeFrequency,
	// 				pickedGenomeFrequency[i],
	// 				pickedGenomeFrequency[i + comparisonJump] ) )
	// 	}
	// }
}

func (s *SelectorSuite) TestShouldRouletteWhenTotalFitnessIs0(t *C) {

	numGenomes := 10
	genomeArray := make([]goga.Genome, numGenomes)
	for i := 0; i < numGenomes; i++ {
		genomeArray[i] = goga.NewGenome(goga.Bitset{})
		genomeArray[i].SetFitness(float64(i))
	}

	goga.Roulette(genomeArray, 0)
}

func (s *SelectorSuite) TestShouldPanicWithMismatchedFitness(t *C) {
	numGenomes := 10
	genomeArray := make([]goga.Genome, numGenomes)
	for i := 0; i < numGenomes; i++ {
		genomeArray[i] = goga.NewGenome(goga.Bitset{})
		genomeArray[i].SetFitness(1)
	}

	// Note: not guaranteed (sp?) to fail, but pretty likely
	t.Assert(func() { goga.Roulette(genomeArray, math.MaxInt32) }, Panics, "total fitness is too large")
}

func (s *SelectorSuite) TestShouldPanicWhenGenomeArrayLengthIs0(t *C) {
	genomeArray := []goga.Genome{}
	t.Assert(len(genomeArray), Equals, 0)
	t.Assert(func() { goga.Roulette(genomeArray, 0) }, Panics, "genome array contains no elements")
}

func (s *SelectorSuite) TestShouldPassBackGenomeFromGenomeArray(t *C) {
	numGenomes := 10
	genomeArray := make([]goga.Genome, numGenomes)

	for i := range genomeArray {
		genomeArray[i] = goga.NewGenome(goga.Bitset{})
		genomeArray[i].SetFitness(1)
	}

	totalFitness := numGenomes
	for i := 0; i < 100; i++ {
		selectedGenome := goga.Roulette(genomeArray, float64(totalFitness))

		found := false
		for i := range genomeArray {
			if genomeArray[i] == selectedGenome {
				found = true
				break
			}
		}

		t.Assert(found, IsTrue)
	}
}

func (s *SelectorSuite) TestShouldConfig_Multiple(t *C) {

	for i := 0; i < 100; i++ {
		numCalls1 := 0
		numCalls2 := 0
		myFunc1 := func(array []goga.Genome, totalFitness float64) goga.Genome {
			numCalls1++
			return array[0]
		}
		myFunc2 := func(array []goga.Genome, totalFitness float64) goga.Genome {
			numCalls2++
			return array[0]
		}

		s := goga.NewSelector(
			[]goga.SelectorFunctionProbability{
				{P: 0.1, F: myFunc1}, // Note probabilities don't add up to 1
				{P: 0.1, F: myFunc2},
			},
		)

		numIterations := 1000
		genomeArray := make([]goga.Genome, 10)
		for i := 0; i < numIterations; i++ {
			s.Go(genomeArray, 100)
		}

		sixtyPercent := (numIterations / 100) * 60
		fourtyPercent := (numIterations / 100) * 40
		t.Assert(numCalls1 < sixtyPercent, IsTrue, Commentf("Num calls [%v] sixty percent [%v]", numCalls1, sixtyPercent))
		t.Assert(numCalls2 < sixtyPercent, IsTrue, Commentf("Num calls [%v] sixty percent [%v]", numCalls2, sixtyPercent))
		t.Assert(numCalls1 > fourtyPercent, IsTrue, Commentf("Num calls [%v] fourty percent [%v]", numCalls1, fourtyPercent))
		t.Assert(numCalls2 > fourtyPercent, IsTrue, Commentf("Num calls [%v] fourty percent [%v]", numCalls2, fourtyPercent))
	}
}
