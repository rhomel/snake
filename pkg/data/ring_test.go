package data

import (
	"testing"
)

func TestRing(t *testing.T) {
	cases := []struct {
		head   int
		offset int
		size   int
		idx    int
	}{
		{
			head:   0,
			offset: 0,
			size:   1,
			idx:    0,
		},
		{
			head:   0,
			offset: 1,
			size:   1,
			idx:    0,
		},
		{
			head:   0,
			offset: 2,
			size:   3,
			idx:    2,
		},
		{
			head:   2,
			offset: 2,
			size:   3,
			idx:    1,
		},
		{
			head:   1,
			offset: 2,
			size:   3,
			idx:    0,
		},
	}
	for _, tc := range cases {
		t.Logf("index(%d, %d, %d) == %d", tc.head, tc.offset, tc.size, tc.idx)
		if got := index(tc.head, tc.offset, tc.size); got != tc.idx {
			t.Errorf("expected %d, received: %d", tc.idx, got)
		}
	}
}

func expectRing(t *testing.T, ring *Ring, want string) {
	t.Helper()
	got := ring.positionsString()
	if want != got {
		t.Errorf("expected order: %s, received: %s", want, got)
	}
}

func TestRingMoveSingle(t *testing.T) {
	ring := &Ring{
		positions: []Position{
			{1, 0},
		},
		head: 0,
		size: 1,
	}

	expectRing(t, ring, "(1,0)")
	ring.Move(Position{2, 0})
	expectRing(t, ring, "(2,0)")
}

func TestRingMoveMultiple(t *testing.T) {
	ring := &Ring{
		positions: []Position{
			{1, 0},
			{2, 0},
			{3, 0},
		},
		head: 0,
		size: 3,
	}

	expectRing(t, ring, "(1,0) (2,0) (3,0)")
	ring.Move(Position{4, 0})
	expectRing(t, ring, "(4,0) (1,0) (2,0)")
	ring.Move(Position{5, 0})
	expectRing(t, ring, "(5,0) (4,0) (1,0)")
	ring.Move(Position{6, 0})
	expectRing(t, ring, "(6,0) (5,0) (4,0)")
}

func TestRingGrowSingle(t *testing.T) {
	const capacity = 5
	ring := NewRing(Position{1, 0}, capacity)
	ring.Grow(Position{2, 0})
	expectRing(t, ring, "(1,0) (2,0)")
}

func TestRingGrowMultiple(t *testing.T) {
	const capacity = 5
	ring := NewRing(Position{1, 0}, capacity)
	ring.Grow(Position{2, 0})
	ring.Grow(Position{3, 0})
	ring.Move(Position{4, 0})
	ring.Move(Position{5, 0})
	expectRing(t, ring, "(5,0) (4,0) (1,0)")

	ring.Grow(Position{6, 0})
	expectRing(t, ring, "(5,0) (4,0) (1,0) (6,0)")
}
