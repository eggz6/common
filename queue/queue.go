package queue

type element struct {
	data interface{}
	used bool
}

type Queue interface {
	Pop() (interface{}, bool)
	Append(val interface{}) bool

	Empty() bool
	Cap() int
	Len() int
	Scan() []interface{}
}

type queue struct {
	read     uint32
	size     uint32
	length   uint32
	write    uint32
	elements []*element
}

func newQueue(size uint32) *queue {
	return &queue{
		read:     0,
		size:     size,
		length:   0,
		write:    0,
		elements: make([]*element, size),
	}
}

func (q *queue) Pop() (interface{}, bool) {
	if q.Empty() {
		return nil, false
	}

	ele := q.elements[q.read]
	res := ele.data

	q.read++
	q.length--

	if q.read >= q.size {
		q.read = 0
	}

	return res, true
}

func (q *queue) Empty() bool {
	return q.length <= 0
}

func (q *queue) Full() bool {
	return q.length >= q.size
}

func (q *queue) Append(val interface{}) bool {
	if q.Full() {
		return false
	}

	ele := q.elements[q.write]
	if ele == nil {
		ele = &element{
			used: false,
		}

		q.elements[q.write] = ele
	}

	q.write++
	q.length++

	if q.write >= q.size {
		q.write = 0
	}

	ele.data = val
	ele.used = true

	return true
}
