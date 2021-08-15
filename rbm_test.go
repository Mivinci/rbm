package rbm

import (
	"reflect"
	"runtime"
	"testing"
)

func assert(t *testing.T, expect, actual interface{}) {
	if !reflect.DeepEqual(expect, actual) {
		_, file, line, _ := runtime.Caller(1)
		t.Fatalf("\nwant %v, but got %v, %s:%d\n", expect, actual, file, line)
	}
}

func TestSearch(t *testing.T) {
	a := []uint16{1, 3, 4, 5, 6, 7, 8}
	b := []uint16{1, 2, 3, 4, 5, 6, 8, 9}
	i, exist := search(b, 8)
	assert(t, true, exist)
	assert(t, 6, i)

	i, exist = search(b, 7)
	assert(t, false, exist)
	assert(t, 6, i)

	i, exist = search(a, 2)
	assert(t, false, exist)
	assert(t, 1, i)

	c := []uint16{1}
	i, exist = search(c, 1)
	assert(t, true, exist)
	assert(t, 0, i)

	var d []uint16
	i, exist = search(d, 1)
	assert(t, false, exist)
	assert(t, 0, i)
}

func TestSplit(t *testing.T) {
	var x uint32 = 0x80fa45e1
	h, l := split(x)
	assert(t, uint16(0x80fa), h)
	assert(t, uint16(0x45e1), l)
}

func TestRBM(t *testing.T) {
	var r RBM
	assert(t, false, r.Get(1))

	r.Add(1)
	assert(t, true, r.Get(1))
	assert(t, false, r.Get(1000000000))

	r.Add(1000000000)
	assert(t, true, r.Get(1000000000))
	assert(t, 2, r.Cardinality())

	r.Remove(1000000000)
	assert(t, false, r.Get(1000000000))
	assert(t, 1, r.Cardinality())
}

func TestRBM_Or(t *testing.T) {
	var r1, r2 RBM
	r1.Add(100)
	r2.Add(200)
	assert(t, false, r1.Get(200))
	r1.Or(&r2)
	assert(t, true, r1.Get(200))
}

func ExampleContainer() {

}
