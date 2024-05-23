package rlzr

import (
	"testing"
)

func monitorExit(t *testing.T) {
	if r := recover(); r != nil {
		t.Errorf("Function panicked with %v", r)
	}
}

func TestInitParams(t *testing.T) {
	defer monitorExit(t)
	InitParams()
}

func TestConstructWritingQueue(t *testing.T) {
	defer monitorExit(t)
	worker_nums := []int{1, 5, 10, 20}
	for _, worker_num := range worker_nums {
		result := ConstructWritingQueue(worker_num)
		switch interface{}(result).(type) {
		case chan packet_metadata:
			// Correct type. Do nothing.
		default:
			t.Errorf("ExampleFunc() = %T; want chan packet_metadata", result)
		}
	}
}
