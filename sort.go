// Package toposort implements in-place topological sorting using both
// BFS (Kahn's Algorithm) and DFS (reverse postorder).
package toposort

import "errors"

// ErrCircular is returned when a cycle is detected in the graph.
var ErrCircular = errors.New("cyclic")

// Vertex represents a node in the graph. It should define Afters(),
// which returns the indices of nodes that depend on this one.
type Vertex interface {
	Afters() []int
}

// BFS performs an in-place topological sort using Kahn's Algorithm.
// If a cycle is detected, it returns ErrCircular.
func BFS[V Vertex](vertices []V, opts ...Option) error {
	n := len(vertices)
	if n < 2 {
		return nil
	}

	cfg := config{}

	for _, opt := range opts {
		opt(&cfg)
	}

	if cfg.bp == nil {
		cfg.bp = &bufferz{}
	}

	inDegree := cfg.bp.IntSlice(n, 0)
	queue := cfg.bp.IntSlice(0, n)
	order := cfg.bp.IntSlice(0, n)

	// Compute in-degrees
	for _, v := range vertices {
		for _, neighbor := range v.Afters() {
			if neighbor >= 0 && neighbor < n {
				inDegree[neighbor]++
			}
		}
	}

	// Collect zero in-degree nodes
	for i := range n {
		if inDegree[i] == 0 {
			queue = append(queue, i)
		}
	}

	// Process queue
	for len(queue) > 0 {
		u := queue[0]
		queue = queue[1:]
		order = append(order, u)

		for _, w := range vertices[u].Afters() {
			if w >= 0 && w < n {
				inDegree[w]--
				if inDegree[w] == 0 {
					queue = append(queue, w)
				}
			}
		}
	}

	// Cycle detected
	if len(order) != n {
		return ErrCircular
	}

	// Apply order in-place
	pos := inDegree
	for i, v := range order {
		pos[v] = i
	}
	for i := range n {
		for pos[i] != i {
			j := pos[i]
			vertices[i], vertices[j] = vertices[j], vertices[i]
			pos[i], pos[j] = pos[j], pos[i]
		}
	}

	return nil
}

// DFS performs an in-place topological sort using a stack-based DFS.
// If a cycle is detected, it returns ErrCircular.
func DFS[V Vertex](vertices []V, opts ...Option) error {
	n := len(vertices)
	if n < 2 {
		return nil
	}

	cfg := config{}

	for _, opt := range opts {
		opt(&cfg)
	}

	if cfg.bp == nil {
		cfg.bp = &bufferz{}
	}

	visited := cfg.bp.BoolSlice(n, 0)
	recStack := cfg.bp.BoolSlice(n, 0)
	order := cfg.bp.IntSlice(0, n)
	stack := cfg.bp.IntSlice(0, n)

	// Iterate through all nodes
	for i := range n {
		if visited[i] {
			continue
		}

		stack = append(stack, i)

		for len(stack) > 0 {
			u := stack[len(stack)-1]

			if !visited[u] {
				visited[u] = true
				recStack[u] = true

				neighbors := vertices[u].Afters()
				neighborSize := len(neighbors)
				for j := range neighborSize {
					w := neighbors[j]
					if w >= 0 && w < n {
						if recStack[w] {
							return ErrCircular
						}
						if !visited[w] {
							stack = append(stack, w)
						}
					}
				}

			} else {
				stack = stack[:len(stack)-1]
				if recStack[u] {
					recStack[u] = false
					order = append(order, u)
				}
			}

			if recStack[u] {
				if len(vertices[u].Afters()) > 0 {
					nextUnvisited := false
					for _, neighbor := range vertices[u].Afters() {
						if neighbor >= 0 && neighbor < n && !visited[neighbor] {
							nextUnvisited = true
							break
						}
					}
					if nextUnvisited {
						continue
					}
				}
				recStack[u] = false
				order = append(order, u)
			}
		}
	}

	// Reverse order for proper sorting
	for i, j := 0, len(order)-1; i < j; i, j = i+1, j-1 {
		order[i], order[j] = order[j], order[i]
	}

	// Apply order in-place
	pos := make([]int, n)
	for i, v := range order {
		pos[v] = i
	}

	for i := range n {
		for pos[i] != i {
			j := pos[i]
			vertices[i], vertices[j] = vertices[j], vertices[i]
			pos[i], pos[j] = pos[j], pos[i]
		}
	}

	return nil
}
