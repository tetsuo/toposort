package toposort_test

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

func ExampleSort() {
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

	_, err := toposort.Sort(map[string]*exampleVertex[string]{
		"Jonas": {id: "Jonas", afters: []string{"Jonas"}},
	})

	fmt.Println(err.Error())

	// Output:
	// [Jonas Sophie Nick Barbara]
	// cyclic: [Jonas Jonas]
}
