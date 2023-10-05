package buffer_go

import (
	"reflect"
	"sync"
	"time"
)

// A Flusher represents logic that need to be done
// before buffer flushed
type Flusher interface {
	// Flush process buffered data
	//
	// As example:
	// accumulate sql inserts and execute as bulk
	Flush([]any) error
}

// A BufferInterface represents interface of buffer
type BufferInterface interface {
	// Push add value to buffer
	Push(v any) error
	// Flush manually flush buffer
	Flush() error
}

// A Buffer is a basic buffer struct with data and buffer configuration
type Buffer struct {
	data       []any
	maxSize    *int
	mu         sync.Mutex
	flushEvery *time.Duration
	flusher    Flusher
}

// Push add value to buffer
func (b *Buffer) Push(v any) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	needFlush := b.maxSize != nil && *b.maxSize <= len(b.data)+1
	if needFlush {
		b.mu.Unlock()
		err := b.Flush()
		if err != nil {
			return nil
		}
		b.mu.Lock()
	}
	b.data = append(b.data, v)
	return nil
}

// Flush manually flush buffer
func (b *Buffer) Flush() error {
	b.mu.Lock()
	defer b.mu.Unlock()
	err := b.flusher.Flush(b.data)
	if err != nil {
		return err
	}
	b.data = nil

	return nil
}

func (b *Buffer) flushByTimer() {
	t := time.NewTimer(*b.flushEvery)
	<-t.C
	_ = b.Flush()
	go b.flushByTimer()
}

// A TypedBuffer is buffer with one value Kind
type TypedBuffer struct {
	Buffer
	// Kind value that will be used to compare pushable value type
	Kind reflect.Kind
}

// Push add value to buffer if value Kind same as TypedBuffer.Kind
func (tb *TypedBuffer) Push(v any) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != tb.Kind {
		return &WrongPushValueType{
			BufferType: tb.Kind.String(),
			ValueType:  rv.Kind().String(),
		}
	}
	return tb.Buffer.Push(v)
}

// NewBuffer created new basic buffer with passed params
//
// flusher: instance of Flusher
//
// maxSize: max size of buffer (if nil, buffer is size-free buffer)
//
// flushEvery: value used to automatically flush the buffer every time at the given
// time value (if nil, the buffer will not be automatically flushed)
func NewBuffer(flusher Flusher, maxSize *int, flushEvery *time.Duration) *Buffer {
	res := Buffer{
		maxSize:    maxSize,
		flushEvery: flushEvery,
		flusher:    flusher,
	}

	if flushEvery != nil {
		go res.flushByTimer()
	}

	return &res
}

// NewTypedBuffer created new typed buffer with passed params
//
// kind: desired kind of buffer item
//
// flusher: instance of Flusher
//
// maxSize: max size of buffer (if nil, buffer is size-free buffer)
//
// flushEvery: value used to automatically flush the buffer every time at the given
// time value (if nil, the buffer will not be automatically flushed)
func NewTypedBuffer(kind reflect.Kind, flusher Flusher, maxSize *int, timeout *time.Duration) *TypedBuffer {
	res := TypedBuffer{
		Buffer: *NewBuffer(flusher, maxSize, timeout),
		Kind:   kind,
	}

	return &res
}
