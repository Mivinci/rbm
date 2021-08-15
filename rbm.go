package rbm

const threshold = 4096

type Container interface {
	Cardinality() int
	Add(v uint16)
	Get(v uint16) bool
	Remove(v uint16)
	Better() Container
	String() string

	// in-place union
	Or(Container) Container
}

type RBM struct {
	ks []uint16
	cs []Container
}

func (r *RBM) Add(x uint32) {
	h, l := split(x)
	r.add(h, l)
}

func (r *RBM) Get(x uint32) bool {
	h, l := split(x)
	return r.get(h, l)
}

func (r *RBM) Remove(x uint32) {
	h, l := split(x)
	r.remove(h, l)
}

func (r *RBM) Optimize() {
	for i := range r.cs {
		r.cs[i] = r.cs[i].Better()
	}
}

func (r *RBM) Cardinality() (car int) {
	for _, c := range r.cs {
		car += int(c.Cardinality())
	}
	return
}

func (r *RBM) add(a, b uint16) {
	i, exist := search(r.ks, a)
	if !exist {
		r.insertKey(i, a)
	}
	c := r.cs[i]
	c.Add(b)
}

func (r *RBM) insertKey(i int, k uint16) {
	r.insert(i, k, &array{})
}

func (r *RBM) insert(i int, k uint16, c Container) {
	r.ks = append(r.ks, 0)
	r.cs = append(r.cs, nil)

	copy(r.ks[i+1:], r.ks[i:])
	copy(r.cs[i+1:], r.cs[i:])

	r.ks[i] = k
	r.cs[i] = c
}

func (r *RBM) get(a, b uint16) bool {
	if i, exist := search(r.ks, a); exist {
		return r.cs[i].Get(b)
	}
	return false
}

func (r *RBM) remove(a, b uint16) {
	if i, exist := search(r.ks, a); exist {
		r.cs[i].Remove(b)
	}
}

func search(a []uint16, x uint16) (int, bool) {
	return searchBetween(a, 0, len(a)-1, x)
}

func searchBetween(a []uint16, i, j int, x uint16) (int, bool) {
	var m int
	for i <= j {
		m = i + (j-i)/2
		if a[m] < x {
			i = m + 1
		} else if a[m] > x {
			j = m - 1
		} else {
			return m, true
		}
	}
	if i <= j && j >= 0 && x > a[i] {
		i++
	}
	return i, false
}

func split(x uint32) (uint16, uint16) {
	h := x & 0xffff0000 >> 16
	l := x & 0x0000ffff
	return uint16(h), uint16(l)
}

func (r *RBM) Or(b *RBM) {
	var i, j int
	n, m := len(r.ks), len(b.ks)
	for i < n && j < m {
		k1 := r.ks[i]
		k2 := b.ks[j]
		if k1 < k2 {
			i++
		} else if k1 > k2 {
			r.insert(i, k2, b.cs[j])
			i++
			j++
		} else {
			r.cs[i] = r.cs[i].Or(b.cs[j])
			i++
			j++
		}
	}
	if j < n {
		r.ks = append(r.ks, b.ks[j:]...)
		r.cs = append(r.cs, b.cs[j:]...)
	}
}
