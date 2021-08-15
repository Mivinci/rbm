package rbm

import "math/bits"

type bitset struct {
	car int // 1-bit count
	val []uint64
}

func (b *bitset) Add(v uint16) {
	n := v / 64
	size := len(b.val)
	if uint16(size) <= n {
		b.grow(n)
	}
	b.val[n] |= 1 << (v % 64)
	b.car++
}

func (b *bitset) AddMany(vs ...uint16) {
	edge := vs[len(vs)-1] / 64
	size := len(b.val)
	if uint16(size) <= edge {
		b.grow(edge)
	}
	for _, v := range vs {
		n := v / 64
		b.val[n] |= 1 << (v % 64)
	}
	b.car += len(vs)
}

func (b *bitset) grow(n uint16) {
	size := len(b.val)
	rb := make([]uint64, n-uint16(size)+1)
	b.val = append(b.val, rb...)
	rb = nil
}

func (b *bitset) Get(v uint16) bool {
	n := v / 64
	size := len(b.val)
	if uint16(size) <= n {
		return false
	}
	return b.val[n]&(1<<(v%64)) != 0
}

func (b *bitset) Remove(v uint16) {
	n := v / 64
	size := len(b.val)
	if uint16(size) < n {
		return
	}
	b.val[n] &^= 1 << (v % 64)
	b.car--
}

// algorithm2 in paper "Better bitmap performance with Roaring bitmaps"
func (b *bitset) ToArray() *array {
	var s []uint16
	for i, w := range b.val {
		for w != 0 {
			t := w & -w
			s = append(s, uint16(i*64+bits.OnesCount64(t-1)))
			w ^= t
		}
	}
	return &array{val: s, car: len(s)}
}

func (b *bitset) Cardinality() int {
	return b.car
}

func (b *bitset) Better() Container {
	if b.car > threshold {
		return b
	}
	return b.ToArray()
}

func (b *bitset) clone() *bitset {
	var c bitset
	c.val = make([]uint64, len(b.val))
	c.car = b.car
	copy(c.val, b.val[:])
	return &c
}

func (b *bitset) Or(c Container) Container {
	switch a := c.(type) {
	case *bitset:
		return b.or(a)
	case *array:
		return a.orBitset(b)
	}
	return nil
}

func (b *bitset) or(a *bitset) Container {
	var c bitset
	var i int
	k := min(len(a.val), len(b.val))
	for i < k {
		c.val = append(c.val, b.val[i]|a.val[i])
		c.car += bits.OnesCount64(c.val[i])
		i++
	}
	if len(a.val) > k {
		c.val = append(c.val, a.val[i:]...)
	}
	if len(b.val) > k {
		c.val = append(c.val, b.val[i:]...)
	}
	return &c
}

func (b *bitset) String() string {
	return b.ToArray().String()
}

func min(a, b int) int {
	if a > b {
		return b
	}
	return a
}
