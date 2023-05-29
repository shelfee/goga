package goga

// IGenome associates a fitness with a bitset
type Genome interface {
	GetFitness() float64
	SetFitness(float64)
	GetBits() *Bitset
	GetOrigin() float64
	SetOrigin(float64)
}

type genome struct {
	fitness float64
	bitset  Bitset
	origin  float64
}

// NewGenome creates a genome with a bitset and
// a zero'd fitness score
func NewGenome(bitset Bitset) Genome {
	return &genome{bitset: bitset}
}

func (g *genome) GetFitness() float64 {
	return g.fitness
}

func (g *genome) SetFitness(fitness float64) {
	g.fitness = fitness
}

func (g *genome) GetBits() *Bitset {
	return &g.bitset
}

func (g *genome) SetOrigin(origin float64) {
	g.origin = origin
}

func (g *genome) GetOrigin() float64 {
	return g.origin
}
