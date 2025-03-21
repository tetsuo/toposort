package toposort

type Buffers interface {
	IntSlice(int, int) []int
	BoolSlice(int, int) []bool
}

type bufferz struct{}

func (b *bufferz) IntSlice(size, capacity int) []int {
	if capacity == 0 {
		return make([]int, size)
	}
	return make([]int, size, capacity)
}

func (b *bufferz) BoolSlice(size, capacity int) []bool {
	if capacity == 0 {
		return make([]bool, size)
	}
	return make([]bool, size, capacity)
}
