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

func main() {
    // Define the relationships in the graph
    relations := map[string]string{
        "Barbara": "Nick",
        "Nick":    "Sophie",
        "Sophie":  "Jonas",
    }

    // Perform topological sorting
    sorted, err := toposort.Sort(relations)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }

    fmt.Println("Sorted order:", sorted)
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
relations := map[string]string{
    "Jonas": "Jonas",
}

_, err := toposort.Sort(relations)
if err != nil {
    fmt.Println("Error:", err)
}
```

**Output:**

```
Error: cyclic: [Jonas Jonas]
```

This indicates that a cycle exists in the graph, making topological sorting impossible.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

