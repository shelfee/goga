package goga

import (
	"sort"
	"sync"
	"time"
)

// GeneticAlgorithm -
// The main component of goga, holds onto the state of the algorithm -
// * Mater - combining evolved genomes
// * EliteConsumer - an optional class that accepts the 'elite' of each population generation
// * Simulator - a simulation component used to score each genome in each generation
// * BitsetCreate - used to create the initial population of genomes
type GeneticAlgorithm struct {
	Mater         Mater
	EliteConsumer EliteConsumer
	Simulator     Simulator
	Selector      Selector
	BitsetCreate  BitsetCreate

	populationSize          int
	MaterExtraRatio         int
	population              []Genome
	totalFitness            float64
	genomeSimulationChannel chan Genome
	exitFunc                func(Genome) bool
	waitGroup               *sync.WaitGroup
	parallelSimulations     int
}

type Options struct {
	PopulationSize      int
	MaterExtraRatio     int
	ParallelSimulations int
}
type Option func(*Options)

func PopulationSize(n int) Option {
	return func(o *Options) {
		o.PopulationSize = n
	}
}

func MaterExtraRatio(n int) Option {
	return func(o *Options) {
		o.MaterExtraRatio = n
	}
}

func ParallelSimulations(n int) Option {
	return func(o *Options) {
		o.ParallelSimulations = n
	}
}

// NewGeneticAlgorithm returns a new GeneticAlgorithm structure with null implementations of
// EliteConsumer, Mater, Simulator, Selector and BitsetCreate
func NewGeneticAlgorithm() GeneticAlgorithm {
	return GeneticAlgorithm{
		EliteConsumer: &NullEliteConsumer{},
		Mater:         &NullMater{},
		Simulator:     &NullSimulator{},
		Selector:      &NullSelector{},
		BitsetCreate:  &NullBitsetCreate{},
	}
}

func (ga *GeneticAlgorithm) createPopulation() []Genome {
	ret := make([]Genome, ga.populationSize)
	for i := 0; i < ga.populationSize; i++ {
		ret[i] = NewGenome(ga.BitsetCreate.Go())
	}
	return ret
}

// Init initialises internal components, sets up the population size
// and number of parallel simulations
func (ga *GeneticAlgorithm) Init(opt ...Option) {
	opts := Options{
		PopulationSize:      10,
		MaterExtraRatio:     2,
		ParallelSimulations: 1,
	}
	for _, o := range opt {
		o(&opts)
	}

	ga.populationSize = opts.PopulationSize
	ga.population = ga.createPopulation()
	ga.parallelSimulations = opts.ParallelSimulations
	ga.MaterExtraRatio = opts.MaterExtraRatio
	ga.waitGroup = new(sync.WaitGroup)
}

func (ga *GeneticAlgorithm) beginSimulation() {
	ga.Simulator.OnBeginSimulation()
	ga.totalFitness = 0
	for i := 0; i < len(ga.population); i++ {
		ga.totalFitness += ga.population[i].GetFitness()
	}
	ga.genomeSimulationChannel = make(chan Genome)

	// todo: make configurable
	for i := 0; i < ga.parallelSimulations; i++ {
		go func(genomeSimulationChannel chan Genome,
			waitGroup *sync.WaitGroup, simulator Simulator) {

			for genome := range genomeSimulationChannel {
				defer waitGroup.Done()
				simulator.Simulate(genome)
			}
		}(ga.genomeSimulationChannel, ga.waitGroup, ga.Simulator)
	}
}

func (ga *GeneticAlgorithm) onNewGenomeToSimulate(g Genome) {
	ga.waitGroup.Add(1)
	ga.genomeSimulationChannel <- g
}

func (ga *GeneticAlgorithm) syncSimulatingGenomes() {
	close(ga.genomeSimulationChannel)
	ga.waitGroup.Wait()
}

func (ga *GeneticAlgorithm) getElite() Genome {
	var ret Genome
	for i := 0; i < ga.populationSize; i++ {
		if ret == nil || ga.population[i].GetFitness() > ret.GetFitness() {
			ret = ga.population[i]
		}
	}
	return ret
}

// SimulateUntil simulates a population until 'exitFunc' returns true
// The 'exitFunc' is passed the elite of each population and should return true
// if the elite reaches a certain criteria (e.g. fitness above a certain threshold)
func (ga *GeneticAlgorithm) SimulateUntil(exitFunc func(Genome) bool) bool {
	ga.exitFunc = exitFunc
	return ga.Simulate()
}

func (ga *GeneticAlgorithm) shouldExit(elite Genome) bool {
	if ga.exitFunc == nil {
		return ga.Simulator.ExitFunc(elite)
	}
	return ga.exitFunc(elite)
}

// Simulate runs the genetic algorithm
func (ga *GeneticAlgorithm) Simulate() bool {

	if ga.populationSize == 0 {
		return false
	}

	ga.beginSimulation()
	for i := 0; i < ga.populationSize; i++ {
		ga.onNewGenomeToSimulate(ga.population[i])
	}
	ga.syncSimulatingGenomes()
	ga.Simulator.OnEndSimulation()

	for {
		elite := ga.getElite()
		ga.Mater.OnElite(elite)
		ga.EliteConsumer.OnElite(elite)
		if ga.shouldExit(elite) {
			break
		}

		time.Sleep(1 * time.Microsecond)

		ga.beginSimulation()

		newPopulationSize := ga.populationSize * ga.MaterExtraRatio
		newPopulation := make([]Genome, newPopulationSize) //ga.createPopulation()
		cache := make(map[string]bool)
		for i := 0; i < newPopulationSize; {
			g1 := ga.Selector.Go(ga.population, ga.totalFitness)
			g2 := ga.Selector.Go(ga.population, ga.totalFitness)
			g3, g4 := ga.Mater.Go(g1, g2)
			k := string(g3.GetBits().GetAll())
			if _, ok := cache[k]; !ok {
				cache[k] = true
				newPopulation[i] = g3
				ga.onNewGenomeToSimulate(newPopulation[i])
				i += 1
			}
			if i < newPopulationSize {
				k = string(g4.GetBits().GetAll())
				if _, ok := cache[k]; !ok {
					cache[k] = true
					newPopulation[i] = g4
					ga.onNewGenomeToSimulate(newPopulation[i])
					i += 1
				}
			}
		}
		ga.syncSimulatingGenomes()
		sort.SliceStable(newPopulation, func(i, j int) bool {
			return newPopulation[i].GetFitness() > newPopulation[j].GetFitness()
		})
		ga.population = newPopulation[:ga.populationSize]
		ga.Simulator.OnEndSimulation()
	}

	return true
}

// GetPopulation returns the population
func (ga *GeneticAlgorithm) GetPopulation() []Genome {
	return ga.population
}
