package rbm

import (
	"strconv"
	"unsafe"
)

type array struct {
	car int
	val []uint16
}

func (a *array) Add(v uint16) {
	i, exist := searchBetween(a.val, 0, a.car-1, v)
	if exist {
		return
	}
	a.val = append(a.val, 0)
	copy(a.val[i+1:], a.val[i:])
	a.val[i] = v
	a.car++
}

func (a *array) Get(v uint16) bool {
	_, exist := searchBetween(a.val, 0, a.car-1, v)
	return exist
}

func (a *array) Remove(v uint16) {
	i, exist := searchBetween(a.val, 0, a.car-1, v)
	if !exist {
		return
	}
	copy(a.val[i:], a.val[i+1:])
	a.car--
}

func (a *array) ToBitset() *bitset {
	var b bitset
	b.AddMany(a.val...)
	return &b
}

func (a *array) Cardinality() int {
	return a.car
}

func (a *array) Better() Container {
	if a.car <= threshold {
		return a
	}
	return a.ToBitset()
}

// func hammingWeight(n uint16) int {
// 	var r int
// 	for n != 0 {
// 		n &= n - 1
// 		r++
// 	}
// 	return r
// }

func (a *array) Or(c Container) Container {
	switch b := c.(type) {
	case *array:
		return a.or(b)
	case *bitset:
		return a.orBitset(b)
	}
	return nil
}

func (a *array) or(b *array) Container {
	var c array
	c.val = merge(a.val, b.val)
	c.car = len(c.val) // new cardinality may be greater than 4096
	return &c
}

func merge(a, b []uint16) (c []uint16) {
	var i, j int
	n, m := len(a), len(b)
	for i < n && j < m {
		if a[i] < b[j] {
			c = append(c, a[i])
			i++
		} else if a[i] > b[j] {
			c = append(c, b[i])
			j++
		} else {
			c = append(c, a[i])
			i++
			j++
		}
	}
	if i < n {
		c = append(c, a[i:]...)
	}
	if j < m {
		c = append(c, b[j:]...)
	}
	return
}

func (a *array) orBitset(b *bitset) Container {
	c := b.clone()
	for _, v := range a.val {
		c.Add(v)
	}
	return c
}

func (a *array) String() string {
	var buf []byte
	buf = append(buf, '{')
	for _, v := range a.val {
		buf = append(buf, strconv.FormatUint(uint64(v), 10)...)
		buf = append(buf, ',')
	}
	buf[len(buf)-1] = '}'
	return *(*string)(unsafe.Pointer(&buf))
}
