package data

import "fmt"

type Position struct {
	x int
	y int
}

type Ring struct {
	positions []Position

	size int
	head int
}

func NewRing(head Position, capacity int) *Ring {
	positions := make([]Position, capacity)
	positions[0] = head
	return &Ring{
		positions: positions,
		head:      0,
		size:      1,
	}
}

func index(head, offset, size int) int {
	return (head + offset) % size
}

func (r *Ring) index(i int) int {
	return index(r.head, i, r.size)
}

func (r *Ring) tail() int {
	return r.index(r.size - 1)
}

func (r *Ring) GetHead() Position {
	return r.positions[r.head]
}

func (r *Ring) GetTail() Position {
	return r.positions[r.tail()]
}

func (r *Ring) GetSize() int {
	return r.size
}

func (r *Ring) Get(i int) Position {
	return r.positions[r.index(i)]
}

func (r *Ring) positionsString() string {
	s := ""
	for i := 0; i < r.size; i++ {
		p := r.Get(i)
		if s != "" {
			s += " "
		}
		s += fmt.Sprintf("(%d,%d)", p.x, p.y)
	}
	return s
}

// move makes position p the new head and shifts all segments by one removing
// the last position. The last position is returned.
func (r *Ring) Move(p Position) Position {
	last := r.GetTail()
	tailIndex := r.tail()
	r.positions[tailIndex] = p
	r.head = tailIndex
	return last
}

// initial      abc
// [a][b][c]
//  h     t
//
// move d ->   dab
// [a][b][d]
//     t  h
//
// move e ->  eda
// [a][e][d]
//  t  h
//
// move f -> fed
// [f][e][d]
//  h     t

func (r *Ring) Grow(p Position) {
	tail := r.tail()
	if tail == r.size-1 {
		r.positions[r.size] = p
		r.size++
	} else {
		// move everything from the head over by 1
		for i := r.size; i > r.head; i-- {
			r.positions[i] = r.positions[i-1]
		}
		// insert the new item at where head is
		r.positions[r.head] = p
		// move head to the proper position
		r.head++
		r.size++
	}
}

func (r *Ring) IsHeadOnBody() bool {
	tail := r.tail()
	for i := r.head; i != tail; i = index(r.head, i+1, r.size) {
		if r.positions[r.head].x == r.positions[i].x &&
			r.positions[r.head].y == r.positions[i].y {
			return true
		}
	}
	return false
}
