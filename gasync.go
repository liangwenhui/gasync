package gasync

import (
	"errors"
	"runtime/debug"
)

func GoAsync[T any](fn func() (T, error)) Feature[T] {
	feature := newDataFeature[T]()
	go func() {
		defer func() {
			switch err := recover().(type) {
			case nil:
				// undo
			case error:
				stack := debug.Stack()
				feature.Done(*new(T), errors.Join(errors.New(string(stack)), err))
			}
		}()
		feature.Done(fn())
	}()
	return feature
}

type doneAble interface {
	IsDone() bool
}

func AllDone(ables ...doneAble) {
	for _, able := range ables {
		able.IsDone()
	}
}
