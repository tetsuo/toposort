# toposort

toposort implements topological sorting in Go using two algorithms: Kahn's Algorithm (BFS) and reverse postorder DFS. It performs in-place sorting and detects cycles.

## Usage

Given edges:

```
5 → 2
5 → 0
4 → 0
4 → 1
2 → 3
3 → 1
```

The graph looks like this:

```
    5       4
   / \     / \
  2   0   0   1
   \
    3
     \
      1
```

### Topological sort order

```go
package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/tetsuo/toposort"
)

type vtx struct {
	id     int
	afters []int
}

func (v vtx) Afters() []int {
	return v.afters
}

func createGraph() []vtx {
	return []vtx{
		{id: 0, afters: nil},
		{id: 1, afters: nil},
		{id: 2, afters: []int{3}},
		{id: 3, afters: []int{1}},
		{id: 4, afters: []int{0, 1}},
		{id: 5, afters: []int{0, 2}},
	}
}

func collect(graph []vtx) string {
	result := make([]string, len(graph))
	for i, v := range graph {
		result[i] = strconv.Itoa(v.id)
	}
	return strings.Join(result, " → ")
}

func main() {
	graph := createGraph()

	if err := toposort.BFS(graph, len(graph)); err != nil {
		if errors.Is(err, toposort.ErrCircular) {
			// cycle detected
		}
		return
	}

	fmt.Println(collect(graph)) // 4 → 5 → 0 → 2 → 3 → 1

	graph = createGraph()
	if err := toposort.DFS(graph, len(graph)); err != nil {
		if errors.Is(err, toposort.ErrCircular) {
			// cycle detected
		}
		return
	}

	fmt.Println(collect(graph)) // 5 → 4 → 2 → 3 → 1 → 0
}
```
