package goga

import (
	"encoding/binary"
	"math"
)

type Float64Requirement struct {
	Precision float64
	MaxValue  float64
	MinValue  float64
	Specific  map[int]struct {
		Precision float64
		MaxValue  float64
		MinValue  float64
	}
}

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

// Round 四舍五入，ROUND_HALF_UP 模式实现
// 返回将 val 根据指定精度 oneUnit（1被分割的单位）进行四舍五入的结果。precision 也可以是负数或零。
// Round(1.7, 0.5) = 1.5
func Round(val float64, oneUnit float64) float64 {
	return math.Floor(val/oneUnit+0.5) * oneUnit
}
