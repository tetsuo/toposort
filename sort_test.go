package toposort_test

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/tetsuo/toposort"
)

type vtx struct {
	id     int
	afters []int
}

func (v vtx) Afters() []int {
	return v.afters
}

func TestSort(t *testing.T) {
	testCases := []struct {
		vs  []vtx
		err error
		ord []string
	}{
		{
			vs: []vtx{
				//    5        4
				//   / \      / \
				//  2   0    0   1
				//   \
				//    3
				//     \
				//      1
				{id: 0, afters: []int{}},
				{id: 1, afters: []int{}},
				{id: 2, afters: []int{3}},
				{id: 3, afters: []int{1}},
				{id: 4, afters: []int{0, 1}},
				{id: 5, afters: []int{0, 2}},
			},
			ord: []string{
				"4,5,0,2,3,1",
				"5,4,2,3,1,0",
			},
		},
		{
			vs: []vtx{
				//    6
				//   / \
				//  3   4
				//   \   \
				//    2   5
				//     \
				//      1
				//       \
				//        0
				{id: 0, afters: []int{}},
				{id: 1, afters: []int{0}},
				{id: 2, afters: []int{1}},
				{id: 3, afters: []int{2}},
				{id: 4, afters: []int{5}},
				{id: 5, afters: []int{}},
				{id: 6, afters: []int{3, 4}},
			},
			ord: []string{
				"6,3,4,2,5,1,0",
				"6,4,5,3,2,1,0",
			},
		},
		{
			vs: []vtx{
				//  6       7
				//  |      / \
				//  3     4   5
				//  |    / \
				//  1   2   0
				//  |
				//  0
				{id: 0, afters: []int{}},
				{id: 1, afters: []int{0}},
				{id: 2, afters: []int{}},
				{id: 3, afters: []int{1}},
				{id: 4, afters: []int{2, 0}},
				{id: 5, afters: []int{}},
				{id: 6, afters: []int{3}},
				{id: 7, afters: []int{4, 5}},
			},
			ord: []string{
				"6,7,3,4,5,1,2,0",
				"7,6,5,4,3,2,1,0",
			},
		},
		{
			vs: []vtx{
				//      7
				//     / \
				//    3   5
				//   / \   \
				//  1   2   6
				//   \      |
				//    0     4
				{id: 0, afters: []int{}},
				{id: 1, afters: []int{0}},
				{id: 2, afters: []int{}},
				{id: 3, afters: []int{1, 2}},
				{id: 4, afters: []int{}},
				{id: 5, afters: []int{6}},
				{id: 6, afters: []int{4}},
				{id: 7, afters: []int{3, 5}},
			},
			ord: []string{
				"7,3,5,1,2,6,0,4",
				"7,5,6,4,3,2,1,0",
			},
		},
		{
			vs: []vtx{
				//    10
				//   /  \
				//  8    9
				//  |   / \
				//  5  6   7
				//  | / \
				//  3    4
				//  |    |
				//  1    2
				//  |
				//  0
				{id: 0, afters: []int{}},
				{id: 1, afters: []int{0}},
				{id: 2, afters: []int{}},
				{id: 3, afters: []int{1}},
				{id: 4, afters: []int{2}},
				{id: 5, afters: []int{3}},
				{id: 6, afters: []int{3, 4}},
				{id: 7, afters: []int{}},
				{id: 8, afters: []int{5}},
				{id: 9, afters: []int{6, 7}},
				{id: 10, afters: []int{8, 9}},
			},
			ord: []string{
				"10,8,9,5,6,7,3,4,1,2,0",
				"10,9,8,7,6,5,4,3,2,1,0",
			},
		},
		{
			vs: []vtx{
				//        9
				//       /
				//      8
				//     /
				//    5       7
				//   / \     / \
				//  1   3   0   2
				//   \ /    |
				//    4     6
				{id: 0, afters: []int{6}},
				{id: 1, afters: []int{4}},
				{id: 2, afters: []int{}},
				{id: 3, afters: []int{4}},
				{id: 4, afters: []int{}},
				{id: 5, afters: []int{1, 3}},
				{id: 6, afters: []int{}},
				{id: 7, afters: []int{0, 2}},
				{id: 8, afters: []int{5}},
				{id: 9, afters: []int{8}},
			},
			ord: []string{
				"7,9,0,2,8,6,5,1,3,4",
				"9,8,7,5,3,2,1,4,0,6",
			},
		},
		{
			vs: []vtx{
				//    8
				//   / \
				//  4   6
				//   \   \
				//    3   7
				//   / \
				//  1   2
				//   \
				//    0
				//    |
				//    5
				{id: 0, afters: []int{5}},
				{id: 1, afters: []int{0}},
				{id: 2, afters: []int{}},
				{id: 3, afters: []int{1, 2}},
				{id: 4, afters: []int{3}},
				{id: 5, afters: []int{}},
				{id: 6, afters: []int{7}},
				{id: 7, afters: []int{}},
				{id: 8, afters: []int{4, 6}},
			},
			ord: []string{
				"8,4,6,3,7,1,2,0,5",
				"8,6,7,4,3,2,1,0,5",
			},
		},
		{
			vs: []vtx{
				//  5 → 3 → 1
				//  ↑    ↓   |
				//  4    2 ←─+
				//  ↓
				//  0
				{id: 0, afters: []int{4}},
				{id: 1, afters: []int{2}},
				{id: 2, afters: []int{3}},
				{id: 3, afters: []int{1}},
				{id: 4, afters: []int{5}},
				{id: 5, afters: []int{3}},
			},
			err: toposort.ErrCircular,
		},
		{
			vs: []vtx{
				//  7 → 4 → 6
				//  ↑    ↓   |
				//  5 ←  3 ← +
				//  ↓
				//  2 → 1 → 0
				//      ↑   ↓
				//      + → 8
				{id: 0, afters: []int{8}},
				{id: 1, afters: []int{0}},
				{id: 2, afters: []int{1}},
				{id: 3, afters: []int{6}},
				{id: 4, afters: []int{3}},
				{id: 5, afters: []int{2, 7}},
				{id: 6, afters: []int{3}},
				{id: 7, afters: []int{4}},
				{id: 8, afters: []int{1}},
			},
			err: toposort.ErrCircular,
		},
		{
			vs: []vtx{},
			ord: []string{
				"",
				"",
			},
		},
	}
	for i, tt := range testCases {
		t.Run(fmt.Sprintf("graph-%d", i), func(t *testing.T) {
			fn := toposort.BFS[vtx]
			for j := range 2 {
				t.Run(fmt.Sprintf("graph-%d-%d", i, j), func(t *testing.T) {
					vs := slices.Clone(tt.vs)
					err := fn(vs)
					if err != nil {
						if tt.err != nil {
							if !errors.Is(err, tt.err) {
								t.Fatalf("expected error %v, got %v", tt.err, err)
							}
						} else {
							t.Fatalf("got unexpected error %v", err)
						}
					} else {
						var ids []string
						for _, v := range vs {
							ids = append(ids, strconv.Itoa(v.id))
						}
						got := strings.Join(ids, ",")
						if got != tt.ord[j] {
							t.Fatalf("expected ordering to be %v got %v", tt.ord[j], got)
						}
					}
					fn = toposort.DFS[vtx]
				})
			}
		})
	}
}
