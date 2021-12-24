package goga_test

import (
	"github.com/tomcraven/goga"
	. "gopkg.in/check.v1"
)

type SimulatorTestSuite struct {
}

func (s *SimulatorTestSuite) SetUpTest(t *C) {
}
func (s *SimulatorTestSuite) TearDownTest(t *C) {
}

var _ = Suite(&SimulatorTestSuite{})

func (s *SimulatorTestSuite) TestShouldReturnFalseFromExitFunctionFromNullSimulator(t *C) {
	nullSimulator := goga.NullSimulator{}
	genome := goga.NewGenome(goga.Bitset{})
	t.Assert(nullSimulator.ExitFunc(genome), IsFalse)
}
