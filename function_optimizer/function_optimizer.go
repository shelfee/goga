package function_optimizer

import (
	"fmt"
	"github.com/tomcraven/goga"
	"math"
	"math/rand"
	"runtime"
)

type funcMaterSimulator struct {
	function       func([]float64) float64
	transFunc      func(float642 float64) float64
	paramSize      int
	minIter        int
	stableExitIter int
	stableMinIter  int
	onStable       func()
	iter           int
	stableIter     int
	lastFitness    float64
	onBegin        func()
	onEnd          func()
}

func (sms *funcMaterSimulator) OnBeginSimulation() {
	if sms.onBegin != nil {
		sms.onBegin()
	}
}
func (sms *funcMaterSimulator) OnEndSimulation() {
	if sms.onEnd != nil {
		sms.onEnd()
	}
}

func (sms *funcMaterSimulator) Simulate(g goga.Genome) {
	params := goga.ParseBitsToFloat64Arr(g.GetBits())
	g.SetOrigin(sms.function(params))
	g.SetFitness(sms.transFunc(g.GetOrigin()))
}
func (sms *funcMaterSimulator) ExitFunc(g goga.Genome) bool {
	if g.GetFitness() != sms.lastFitness {
		sms.stableIter = 0
	}
	sms.lastFitness = g.GetFitness()
	sms.stableIter += 1
	sms.iter += 1
	if sms.minIter < sms.iter && sms.stableIter < sms.stableExitIter {
		return true
	}
	if sms.stableIter > sms.stableMinIter {
		if sms.onStable != nil {
			sms.onStable()
		}
	}
	return false
}

type myBitsetCreate struct {
	paramsSize  int
	initValue   []float64
	requirement *goga.Float64Requirement
}

func (bc *myBitsetCreate) Go() goga.Bitset {
	b := goga.Bitset{}
	b.Create(bc.paramsSize * 8)
	for i := 0; i < bc.paramsSize; i++ {
		var param float64
		if require, ok := bc.requirement.Specific[i]; ok {
			if i < len(bc.initValue) && bc.initValue[i] >= require.MinValue && bc.initValue[i] <= require.MaxValue {
				param = bc.initValue[i]
			} else {
				param = goga.Round(rand.Float64()*(require.MaxValue-require.MinValue)+require.MinValue, require.Precision)
			}
		} else {
			if i < len(bc.initValue) && bc.initValue[i] >= bc.requirement.MinValue && bc.initValue[i] <= bc.requirement.MaxValue {
				param = bc.initValue[i]
			} else {
				param = goga.Round(rand.Float64()*(bc.requirement.MaxValue-bc.requirement.MinValue)+bc.requirement.MinValue, bc.requirement.Precision)
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
	onElite     func(g goga.Genome)
}

func (ec *myEliteConsumer) OnElite(g goga.Genome) {
	if ec.onElite != nil {
		ec.onElite(g)
		return
	}
	gBits := g.GetBits()
	ec.currentIter++
	params := goga.ParseBitsToFloat64Arr(gBits)
	fmt.Println(ec.currentIter, "\t", params, "\tfunc value: ", g.GetOrigin(), "\tfitness: ", g.GetFitness())
}

type Options struct {
	requirement     *goga.Float64Requirement
	paramSize       int
	function        func([]float64) float64
	transFunc       func(float64) float64
	minIter         int
	stableExitIter  int
	stableMinIter   int
	populationSize  int
	numThreads      int
	materExtraRatio int
	lruSize         int
	randomRatio     float64
	onBegin         func()
	onEnd           func()
	onElite         func(g goga.Genome)
	onStable        func()

	mater    goga.Mater
	selector goga.Selector
}
type Option func(*Options)

func OnBegin(n func()) Option {
	return func(o *Options) {
		o.onBegin = n
	}
}
func OnEnd(n func()) Option {
	return func(o *Options) {
		o.onEnd = n
	}
}
func OnElite(n func(g goga.Genome)) Option {
	return func(o *Options) {
		o.onElite = n
	}
}
func OnStable(n func()) Option {
	return func(o *Options) {
		o.onStable = n
	}
}
func PopulationSize(n int) Option {
	return func(o *Options) {
		o.populationSize = n
	}
}
func LRUSize(n int) Option {
	return func(o *Options) {
		o.lruSize = n
	}
}
func NumThreads(n int) Option {
	return func(o *Options) {
		o.numThreads = n
	}
}

func RandomRatio(n float64) Option {
	return func(o *Options) {
		o.randomRatio = n
	}
}
func MaterExtraRatio(n int) Option {
	return func(o *Options) {
		o.materExtraRatio = n
	}
}
func Mater(n goga.Mater) Option {
	return func(o *Options) {
		o.mater = n
	}
}
func Selector(n goga.Selector) Option {
	return func(o *Options) {
		o.selector = n
	}
}
func StableExitIter(n int) Option {
	return func(o *Options) {
		o.stableExitIter = n
	}
}

func StableMinIter(n int) Option {
	return func(o *Options) {
		o.stableMinIter = n
	}
}
func MinIter(n int) Option {
	return func(o *Options) {
		o.minIter = n
	}
}
func Requirement(n *goga.Float64Requirement) Option {
	return func(o *Options) {
		o.requirement = n
	}
}
func TransFunc(n func(float64) float64) Option {
	return func(o *Options) {
		o.transFunc = n
	}
}
func Function(n func([]float64) float64) Option {
	return func(o *Options) {
		o.function = n
	}
}
func ParamSize(n int) Option {
	return func(o *Options) {
		o.paramSize = n
	}
}

func NewFuncAlgo(o ...Option) goga.GeneticAlgorithm {
	opts := Options{
		requirement: &goga.Float64Requirement{
			Precision: 1,
			MinValue:  -math.MaxFloat64,
			MaxValue:  math.MaxFloat64,
		},
		paramSize:       0,
		function:        func(float64s []float64) float64 { return rand.Float64() },
		transFunc:       func(f float64) float64 { return f },
		minIter:         200,
		stableExitIter:  50,
		stableMinIter:   10,
		populationSize:  600,
		lruSize:         2400,
		numThreads:      runtime.NumCPU() - 1,
		materExtraRatio: 4,
		randomRatio:     0.3,
	}
	for _, o := range o {
		o(&opts)
	}
	mater := goga.FloatMater{
		Float64Requirement: *opts.requirement,
	}
	opts.mater = goga.NewMater(
		[]goga.MaterFunctionProbability{
			{P: 0.1, F: mater.ArithmeticMutate},
			{P: 1.0, F: mater.ArithmeticExchange},
			{P: 1.0, F: mater.ArithmeticCrossover, UseElite: true},
		})
	opts.selector = goga.NewSelector(
		[]goga.SelectorFunctionProbability{
			{P: 0.1, F: goga.RandomSelect},
			{P: 1.0, F: goga.Roulette},
		},
	)
	genAlgo := goga.NewGeneticAlgorithm()
	s := funcMaterSimulator{
		paramSize:      opts.paramSize,
		function:       opts.function,
		transFunc:      opts.transFunc,
		stableExitIter: opts.stableExitIter,
		stableMinIter:  opts.stableMinIter,

		minIter:  opts.minIter,
		onBegin:  opts.onBegin,
		onEnd:    opts.onEnd,
		onStable: opts.onStable,
	}
	genAlgo.Simulator = &s
	genAlgo.BitsetCreate = &myBitsetCreate{paramsSize: opts.paramSize, requirement: opts.requirement}
	genAlgo.EliteConsumer = &myEliteConsumer{
		onElite: opts.onElite,
	}
	genAlgo.Mater = opts.mater
	genAlgo.Selector = opts.selector
	genAlgo.Init(goga.LRUSize(opts.lruSize), goga.PopulationSize(opts.populationSize), goga.ParallelSimulations(opts.numThreads), goga.MaterExtraRatio(opts.materExtraRatio), goga.RandomRatio(opts.randomRatio))
	return genAlgo
}
