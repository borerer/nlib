package utils

import (
	"errors"
	"testing"
)

func TestMustNil(t *testing.T) {
	defer func() {
		r := recover()
		if r != nil {
			t.Fatal("should not panic")
		}
	}()
	Must(nil)
}

func TestMustPanic(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("should panic")
		}
		err, ok := r.(error)
		if !ok {
			t.Fatal("should throw error")
		}
		if err.Error() != "some error" {
			t.Fatal("should throw raw error")
		}
	}()
	Must(errors.New("some error"))
}
