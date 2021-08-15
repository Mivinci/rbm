package rbm

import "testing"

func TestBitset(t *testing.T) {
	var b bitset
	assert(t, false, b.Get(1))
	b.Add(1)
	assert(t, true, b.Get(1))
	assert(t, uint16(1), b.Cardinality())

	b.Add(10000)
	assert(t, true, b.Get(10000))
	assert(t, uint16(2), b.Cardinality())

	b.Remove(10000)
	b.Remove(1)
	assert(t, false, b.Get(10000))
	assert(t, false, b.Get(10000))
	assert(t, uint16(0), b.Cardinality())

	b.AddMany(100, 200, 300)
	assert(t, true, b.Get(100))
	assert(t, true, b.Get(200))
	assert(t, true, b.Get(300))
	assert(t, uint16(3), b.Cardinality())
}

func TestBitset_Or(t *testing.T) {
	var a array
	var b bitset
	a.Add(400)
	b.AddMany(100, 200, 300)
	c := b.Or(&a)
	assert(t, "{100,200,300,400}", c.String())

	var b1 bitset
	b1.Add(500)
	c = c.Or(&b1)
	assert(t, "{100,200,300,400,500}", c.String())
}

func TestBitset_ToArray(t *testing.T) {
	var b bitset
	b.AddMany(100, 200, 300)
	a := b.ToArray()
	assert(t, true, a.Get(100))
	assert(t, true, a.Get(200))
	assert(t, true, a.Get(300))
	assert(t, false, a.Get(400))
}
