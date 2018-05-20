# Weight

## Getting start

You can use this package by go get.

```shell
go get github.com/hackez/weight
```

### Smooth weighted round-robin balancing

Using Smooth weighted round-robin balancing design by `nginx`:

```go
package main

import (
    "fmt"

    "github.com/hackez/weight"
)

func main() {
    sw := &weight.SW{}
    sw.Add("A", 4)
    sw.Add("B", 2)
    sw.Add("C", 1)

    total := sw.Total()

    for i := 0; i < total; i++ {
        elem := sw.Next()
        fmt.Printf("%v ", elem)
    }

    // output: A B A C A B A
}
```

### LVS weight round-robin balancing

[Job Scheduling Algorithms in Linux Virtual Server](http://www.linuxvirtualserver.org/docs/scheduling.html)

Using LVS balancing:

```go
package main

import (
    "fmt"

    "github.com/hackez/weight"
)

func main() {
    // TODO:
}
```
