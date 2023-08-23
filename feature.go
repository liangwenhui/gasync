package gasync

import (
	"sync/atomic"
	"time"
)

// Feature : a feature returns when call the GoAsync
type Feature[T any] interface {
	// Get : Keep blocking until complete or abnormal
	Get() (T, error)
	// TryGet : Based or Get add timeout
	TryGet(duration time.Duration) (T, error)
	// GetNow : Not-blocking
	GetNow() (T, bool, error)
	// GetSuccess : Keep blocking until complete or abnormal, but not return errors
	GetSuccess() T
	IsDone() bool
}

var _ Feature[struct{}] = &DataFeature[struct{}]{}

const (
	Running = iota
	Donning
	Done
)

type DataFeature[T any] struct {
	status int32
	sc     chan int
	data   T
	err    error
}

func newDataFeature[T any]() *DataFeature[T] {
	return &DataFeature[T]{
		sc: make(chan int, 0),
	}
}

func (d *DataFeature[T]) Get() (T, error) {
	if !d.IsDone() {
		d.await()
	}
	return d.data, d.err
}

func (d *DataFeature[T]) TryGet(duration time.Duration) (T, error) {
	if !d.IsDone() {
		d.awaitTimeout(duration)
	}
	now, _, err := d.GetNow()
	return now, err
}

func (d *DataFeature[T]) GetNow() (T, bool, error) {
	if !d.IsDone() {
		return *new(T), false, nil
	}
	return d.data, true, d.err
}

func (d *DataFeature[T]) GetSuccess() T {
	get, _ := d.Get()
	return get
}

func (d *DataFeature[T]) IsDone() bool {
	return d.status == Done
}

func (d *DataFeature[T]) await() {
	<-d.sc
}
func (d *DataFeature[T]) awaitTimeout(expire time.Duration) {
	if expire < time.Second*0 {
		return
	}
	select {
	case <-d.sc:
	case <-time.After(expire):
	}
}
func (d *DataFeature[T]) updateStatus(old int32, new int32) bool {
	return atomic.CompareAndSwapInt32(&d.status, old, new)
}

func (d *DataFeature[T]) Done(data T, err error) {
	if !d.updateStatus(Running, Donning) {
		//errors.Join(errors.New(string(stack)), err)
		return
	}
	d.data = data
	d.err = err
	if !d.updateStatus(Donning, Done) {
		//	return errors.New("illegal status")
		d.data = *new(T)
	}
	close(d.sc)
}
