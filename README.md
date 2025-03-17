# toposort

**toposort** performs topological sorting on directed acyclic graphs (DAGs). Topological sorting is the linear ordering of vertices such that for every directed edge `U → V`, vertex `U` comes before vertex `V` in the ordering.

## Installation

To install the `toposort` package, use the following command:

```sh
go get github.com/tetsuo/toposort
```

## Usage

```go
package main

import (
	"fmt"

	"github.com/tetsuo/toposort"
)

type exampleVertex[K comparable] struct {
	afters []K
	id     K
}

func (v *exampleVertex[K]) Afters() []K {
	return v.afters
}

func (v *exampleVertex[K]) ID() K {
	return v.id
}

func main() {
	relations := map[string]string{
		"Barbara": "Nick",
		"Nick":    "Sophie",
		"Sophie":  "Jonas",
	}

	vertices := make(map[string]*exampleVertex[string])

	for c, p := range relations {
		if _, ok := vertices[c]; !ok {
			vertices[c] = &exampleVertex[string]{id: c}
		}
		var e *exampleVertex[string]
		if _, ok := vertices[p]; !ok {
			e = &exampleVertex[string]{id: p}
			vertices[p] = e
		} else {
			e = vertices[p]
		}
		e.afters = append(e.afters, c)
	}

	sorted, _ := toposort.Sort(vertices)

	fmt.Println(sorted)
}
```

**Output:**

```
Sorted order: [Jonas Sophie Nick Barbara]
```

In this example, the `relations` map defines a set of dependencies where each key depends on its corresponding value. The `Sort` function processes these relationships and returns a slice of strings representing the topologically sorted order.

## Error Handling

If the graph contains cycles, the `Sort` function will return an error indicating the presence of a cycle. For example:

```go
_, err := toposort.Sort(map[string]*exampleVertex[string]{
    "Jonas": {id: "Jonas", afters: []string{"Jonas"}},
})

fmt.Println(err.Error())
```

**Output:**

```
Error: cyclic: [Jonas Jonas]
```

This indicates that a cycle exists in the graph, making topological sorting impossible.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

