package toposort

import (
	"errors"
	"fmt"
	"slices"
)

var (
	// ErrCircular is raised when a cyclic relationship has been found.
	ErrCircular = errors.New("cyclic")
	// ErrMultipleRoots is raised when a graph contains multiple root nodes.
	ErrMultipleRoots = errors.New("multiple roots")
)

type Vertex[K comparable] interface {
	Afters() []K
	ID() K
}

// tsort sorts the given graph topologically.
func tsort[K comparable, V Vertex[K]](g map[K]V) (sorted []K, recursive map[K]bool, recursion []K) {
	sorted = []K{}
	visited := make(map[K]bool)
	recursive = make(map[K]bool) // keys caught in a recursive chain
	recursion = []K{}            // recursion paths for printing out in the error messages

	var visit func(id K, ancestors []K)

	visit = func(id K, ancestors []K) {
		vertex := g[id]
		if _, ok := visited[id]; ok {
			return
		}
		ancestors = append(ancestors, id)
		visited[id] = true
		for _, afterID := range vertex.Afters() {
			if slices.Contains(ancestors, afterID) {
				recursive[id] = true
				for _, id := range ancestors {
					recursive[id] = true
				}
				recursion = append(recursion, append([]K{id}, ancestors...)...)
			} else {
				visit(afterID, ancestors[:])
			}
		}
		sorted = append([]K{id}, sorted...)
	}

	for k := range g {
		visit(k, []K{})
	}

	return
}

type Graph[K comparable, V Vertex[K]] struct {
	data      map[K]V    // graph itself
	sorted    []K        // toposorted keys
	recursive map[K]bool // recursive keys
	recursion []K        // recursion paths
}

func Sort[K comparable, V Vertex[K]](vertices map[K]V) ([]K, error) {
	graph := new(Graph[K, V])
	graph.sorted, graph.recursive, graph.recursion = tsort(vertices)
	graph.data = vertices
	if err := validateGraph(graph); err != nil {
		return nil, err
	}
	return graph.sorted, nil
}

// validateGraph checks a graph for recursive paths and multiple root nodes.
func validateGraph[K comparable, V Vertex[K]](g *Graph[K, V]) (err MultiError) {
	var visit func(id K)

	length := 0

	visit = func(id K) {
		v := g.data[id]
		for _, afterID := range v.Afters() {
			length += 1
			visit(afterID)
		}
	}

	roots := []K{}
	var o int
	for _, id := range g.sorted {
		o = length
		length = 0
		if !g.recursive[id] { // avoid stack overflow
			visit(id)
			if length > o {
				// if the length of the dependencies is increased,
				// that means we are traversing a new tree.
				roots = append(roots, id)
			}
		}
	}

	var zero K

	recursions, start, ko := [][]K{}, 0, zero
	for i, k := range g.recursion {
		if k == ko {
			recursions = append(recursions, g.recursion[start:i+1])
			start = i + 1
			ko = zero
		} else if ko == zero {
			ko = k
		}
	}

	// add all cyclic dependency errors to the multierror instance
	for _, xs := range recursions {
		err = append(err, fmt.Errorf("%w: %v", ErrCircular, xs))
	}

	// add multiple roots error after that if found any
	if len(roots) > 1 {
		err = append(err, fmt.Errorf("%w: %v", ErrMultipleRoots, roots))
	}

	return
}
