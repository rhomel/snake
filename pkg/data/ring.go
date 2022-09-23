package data

import (
	"fmt"
)

type Position struct {
	X int
	Y int
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
		s += fmt.Sprintf("(%d,%d)", p.X, p.Y)
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
	if r.size == 1 {
		return false
	}
	for i := 1; i < r.size; i++ {
		idx := index(r.head, i, r.size)
		if r.positions[r.head].X == r.positions[idx].X &&
			r.positions[r.head].Y == r.positions[idx].Y {
			return true
		}
	}
	return false
}

func (r *Ring) HasPosition(p Position) bool {
	for i := 0; i < r.size; i++ {
		idx := index(r.head, i, r.size)
		if r.positions[p.X].X == r.positions[idx].X &&
			r.positions[p.Y].Y == r.positions[idx].Y {
			return true
		}
	}
	return false
}
