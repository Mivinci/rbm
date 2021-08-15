package rbm

import (
	"reflect"
	"testing"
)

func TestMerge(t *testing.T) {
	testcases := []struct {
		a, b, want []uint16
	}{
		{
			[]uint16{1, 2, 3, 4, 5},
			[]uint16{4, 5, 6, 7},
			[]uint16{1, 2, 3, 4, 5, 6, 7},
		},
		{
			[]uint16{1, 2, 3},
			[]uint16{5, 6, 7},
			[]uint16{1, 2, 3, 5, 6, 7},
		},
		{
			[]uint16{1, 2, 3, 4},
			[]uint16{4, 5, 6, 7},
			[]uint16{1, 2, 3, 4, 5, 6, 7},
		},
	}

	for _, tc := range testcases {
		c := merge(tc.a, tc.b)
		if !reflect.DeepEqual(tc.want, c) {
			t.Fatalf("want %v:%d, but got %v:%d", tc.want, len(tc.want), c, len(c))
		}
	}
}

func TestArray(t *testing.T) {
	var a array
	assert(t, false, a.Get(1))
	a.Add(1)
	assert(t, true, a.Get(1))

	a.Add(10000)
	assert(t, true, a.Get(10000))

	a.Remove(10000)
	a.Remove(1)
	assert(t, false, a.Get(10000))
	assert(t, false, a.Get(10000))
}

func TestArray_Or(t *testing.T) {
	var a array
	var b bitset

	a.Add(100)
	assert(t, "{100}", a.String())
	b.AddMany(100, 200, 300)
	assert(t, "{100,200,300}", b.String())
	c := a.Or(&b)
	assert(t, "{100,200,300}", c.String())

	var a1 array
	a1.Add(400)
	c = c.Or(&a1)
	assert(t, "{100,200,300,400}", c.String())
}

func TestArray_ToBitset(t *testing.T) {
	var a array
	a.Add(100)
	a.Add(200)
	a.Add(300)
	b := a.ToBitset()
	assert(t, true, b.Get(100))
	assert(t, true, b.Get(200))
	assert(t, true, b.Get(300))
	assert(t, false, a.Get(400))
}
