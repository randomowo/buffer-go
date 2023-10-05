package buffer_go

import (
	"fmt"
)

type WrongPushValueType struct {
	BufferType string
	ValueType  string
}

func (e *WrongPushValueType) Error() string {
	return fmt.Sprintf("buffer and push value types differ: %s != %s", e.BufferType, e.ValueType)
}
