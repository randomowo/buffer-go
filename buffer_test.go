package buffer_go

import (
	"reflect"
	"testing"
	"time"
)

type MockedFlusher struct{}

func (m *MockedFlusher) Flush(data []any) error {
	return nil
}

func checkData(data, expected []any) bool {
	for i := range data {
		if data[i] != expected[i] {
			return false
		}
	}
	return true
}

func TestBuffer(t *testing.T) {
	buf := NewBuffer(&MockedFlusher{}, nil, nil)
	_ = buf.Push(1)
	_ = buf.Push(2)

	expected := []any{1, 2}

	if len(buf.data) != len(expected) {
		t.Errorf("buffer size differ from expected: %d != %d", len(buf.data), len(expected))
	}

	if !checkData(buf.data, expected) {
		t.Errorf("buffer differ from expected: %+v %+v", buf.data, expected)
	}
}

func TestBuffer_WithMaxSize(t *testing.T) {
	size := 1
	buf := NewBuffer(&MockedFlusher{}, &size, nil)
	_ = buf.Push(1)
	_ = buf.Push(2)

	expected := []any{2}

	if len(buf.data) != len(expected) {
		t.Errorf("buffer size differ from expected: %d != %d", len(buf.data), len(expected))
	}

	if !checkData(buf.data, expected) {
		t.Errorf("buffer differ from expected: %+v %+v", buf.data, expected)
	}
}

func TestBuffer_WithTimeoutFlush(t *testing.T) {
	timeout := time.Second * 2

	buf := NewBuffer(&MockedFlusher{}, nil, &timeout)
	_ = buf.Push(1)
	_ = buf.Push(2)

	expected := []any{1, 2}

	if len(buf.data) != len(expected) {
		t.Errorf("buffer size differ from expected: %d != %d", len(buf.data), len(expected))
	}

	if !checkData(buf.data, expected) {
		t.Errorf("buffer differ from expected: %+v %+v", buf.data, expected)
	}

	time.Sleep(time.Second * 1)

	if len(buf.data) != len(expected) {
		t.Errorf("buffer size differ from expected: %d != %d", len(buf.data), len(expected))
	}

	if !checkData(buf.data, expected) {
		t.Errorf("buffer differ from expected: %+v %+v", buf.data, expected)
	}

	time.Sleep(time.Second * 2)

	if len(buf.data) > 0 {
		t.Errorf("buffer not flushed after %s", timeout.String())
	}
}

func TestTypedBuffer(t *testing.T) {
	buf := NewTypedBuffer(reflect.Int, &MockedFlusher{}, nil, nil)
	_ = buf.Push(1)
	_ = buf.Push(2)

	expected := []any{1, 2}

	if len(buf.data) != len(expected) {
		t.Errorf("buffer size differ from expected: %d != %d", len(buf.data), len(expected))
	}

	if !checkData(buf.data, expected) {
		t.Errorf("buffer differ from expected: %+v %+v", buf.data, expected)
	}

	err := buf.Push("test")
	if err == nil {
		t.Errorf("typed buffer accepted wrong kind")
	}
}
