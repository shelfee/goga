package goga

// Bitset - a simple bitset implementation
type Bitset struct {
	size int
	bits []byte
}

// Create - creates a bitset of length 'size'
func (b *Bitset) Create(size int) {
	b.size = size
	b.bits = make([]byte, size)
}

// GetSize returns the size of the bitset
func (b *Bitset) GetSize() int {
	return b.size
}

// Get returns the value in the bitset at index 'index'
// or -1 if the index is out of range
func (b *Bitset) Get(index int) int {
	if index < b.size {
		return int(b.bits[index])
	}
	return -1
}

// GetAll returns an int array of all the bits in the bitset
func (b *Bitset) GetAll() []byte {
	return b.bits
}

func (b *Bitset) setImpl(index, value int) {
	b.bits[index] = byte(value)
}

// Set assigns value 'value' to the bit at index 'index'
func (b *Bitset) Set(index, value int) bool {
	if index < b.size {
		b.setImpl(index, value)
		return true
	}
	return false
}

// SetAll assigns the value 'value' to all the bits in the set
func (b *Bitset) SetAll(value int) {
	for i := 0; i < b.size; i++ {
		b.setImpl(i, value)
	}
}

// SetAll assigns the value 'value' to all the bits in the set
func (b *Bitset) SetAllArr(value []byte) {
	b.bits = value
}

// CreateCopy returns a bit for bit copy of the bitset
func (b *Bitset) CreateCopy() Bitset {
	newBitset := Bitset{}
	newBitset.Create(b.size)
	for i := 0; i < b.size; i++ {
		newBitset.Set(i, b.Get(i))
	}
	return newBitset
}

// Slice returns an out of memory slice of the current bitset
// between bits 'startingBit' and 'startingBit + size'
func (b *Bitset) Slice(startingBit, size int) Bitset {
	ret := Bitset{}
	ret.Create(size)
	ret.bits = b.bits[startingBit : startingBit+size]
	return ret
}

func ParseBitsToFloat64Arr(b *Bitset) []float64 {
	bits := b.GetAll()
	params := make([]float64, len(bits)/8)
	for i := 0; i < len(bits)/8; i++ {
		paramBytes := bits[i*8 : i*8+8]
		params[i] = ByteToFloat64(paramBytes)
	}
	return params
}

func ParseFloat64ArrToBits(arr []float64) *Bitset {
	b := Bitset{}
	b.Create(len(arr) * 8)
	for i := 0; i < len(arr); i++ {
		var param = arr[i]
		byteArr := Float64ToByte(param)
		for idx := 0; idx < 8; idx++ {
			b.Set(i*8+idx, int(byteArr[idx]))
		}
	}
	return &b
}
