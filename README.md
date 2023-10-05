# Buffer-go

Buffer structs for any service with some features

## Usage

### Buffer with timeout flush

Example:

```go
package main

import (
	"time"

	"github.com/randomowo/buffer-go"
)

type MyFlusher struct{}

func (f *MyFlusher) Flush(data []any) error {
	return nil
}

func main() {
	timeout := time.Second
	flusher := new(MyFlusher)
	buffer := buffer_go.NewBuffer(flusher, nil, &timeout)
	_ = buffer.Push(1) // [1]
	_ = buffer.Push(2) // [1 2]
	time.Sleep(time.Second * 2)
	_ = buffer.Push(3) // [3]
}
```

### Sized buffer

Example:

```go
package main

import (
	"github.com/randomowo/buffer-go"
)

type MyFlusher struct{}

func (f *MyFlusher) Flush(data []any) error {
	return nil
}

func main() {
	size := 2
	flusher := new(MyFlusher)
	buffer := buffer_go.NewBuffer(flusher, &size, nil)
	_ = buffer.Push(1) // [1]
	_ = buffer.Push(2) // [1 2]
	_ = buffer.Push(3) // [3]
}
```

### Size-free buffer

Example:

```go
package main

import (
	"github.com/randomowo/buffer-go"
)

type MyFlusher struct{}

func (f *MyFlusher) Flush(data []any) error {
	return nil
}

func main() {
	flusher := new(MyFlusher)
	buffer := buffer_go.NewBuffer(flusher, nil, nil)
	_ = buffer.Push(1) // [1]
	_ = buffer.Push(2) // [1 2]
	_ = buffer.Push(3) // [1 2 3]
	_ = buffer.Flush() // manually flush buffer
	_ = buffer.Push(4) // [4]
}
```