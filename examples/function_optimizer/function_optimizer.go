package main

import (
	"encoding/binary"
	"fmt"
	"github.com/tomcraven/goga"
	"math"
	"math/rand"
)

// Float64ToByte Float64转byte
func Float64ToByte(float float64) []byte {
	bits := math.Float64bits(float)
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, bits)
	return bytes
}

// ByteToFloat64 byte转Float64
func ByteToFloat64(bytes []byte) float64 {
	bits := binary.LittleEndian.Uint64(bytes)
	return math.Float64frombits(bits)
}

type funcMaterSimulator struct {
	function      func([]float64) float64
	paramSize     int
	findSmallest  bool
	transEndValue float64
	transMaxValue float64
	powValue      float64
}

func (sms *funcMaterSimulator) TransToPos(v float64) float64 {
	if v < sms.transEndValue {
		return 0
	}
	return sms.transMaxValue / (1 + math.Pow(sms.powValue, -v))
}
func (sms *funcMaterSimulator) ParseBits(b *goga.Bitset) []float64 {
	bits := b.GetAll()
	params := make([]float64, sms.paramSize)
	for i := 0; i < sms.paramSize; i++ {
		paramBytes := bits[i*8 : i*8+8]
		params[i] = ByteToFloat64(paramBytes)
	}
	return params
}

func (sms *funcMaterSimulator) OnBeginSimulation() {
}
func (sms *funcMaterSimulator) OnEndSimulation() {
}

func (sms *funcMaterSimulator) Simulate(g goga.Genome) {
	params := sms.ParseBits(g.GetBits())
	flag := 1.
	if sms.findSmallest {
		flag = -1.
	}
	g.SetFitness(sms.TransToPos(flag * sms.function(params)))
}
func (sms *funcMaterSimulator) ExitFunc(g goga.Genome) bool {
	return g.GetFitness() > sms.transMaxValue*0.99999999
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
			param = rand.Float64()
		}
		byteArr := Float64ToByte(param)
		for idx := 0; idx < 8; idx++ {
			b.Set(i*8+idx, int(byteArr[idx]))
		}
	}
	return b
}

func main() {
	s := funcMaterSimulator{
		paramSize: 4,
		function: func(float64s []float64) float64 {
			return -math.Pow(float64s[0], 4) + math.Pow(float64s[1], 3) + math.Pow(float64s[2], 2) + math.Pow(float64s[3], 1) + 1
		},
		findSmallest:  false,
		transEndValue: -500,
		transMaxValue: 1000,
		powValue:      1.001,
	}
	b := myBitsetCreate{
		paramsSize: 4,
		initValue:  []float64{2., 4., 1., 3.},
	}
	bitset := b.Go()
	params := s.ParseBits(&bitset)
	for i, p := range params {
		fmt.Printf("idx: %d, param: %f\n", i, p)
	}
}
