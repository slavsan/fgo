package fgo

import (
	"fmt"
	"reflect"
)

func Map[T any, M any](f func(T) M) func(T) M {
	return func(x T) M {
		return f(x)
	}
}

func Filter[T any](f func(T) bool) func(T) bool {
	return func(x T) bool {
		return f(x)
	}
}

func Take(num int) func() int {
	return func() int {
		return num
	}
}

func Pipe[T any, M any, R any](xs []T, fs ...interface{}) ([]M, R, error) {
	var (
		ys    = make([]M, 0, len(xs))
		j     R
		count int
	)
	for i, x := range xs {
		var m M
		if v, ok := any(x).(M); ok {
			m = v
		}

		for _, f := range fs {
			// TODO: on first iteration, check the order of functions is valid, i.e. if reduce is defined, it should be last
			switch f2 := f.(type) {
			case func(T) bool: // filter
				{
					if !f2(x) {
						goto skip
					}
				}
			case func(M) bool: // filter
				{
					if !f2(m) {
						goto skip
					}
				}
			case func(T) T: // map
				{
					x = f2(x)
					if v, ok := any(x).(M); ok {
						m = v
					}
				}
			case func(T) M: // map
				{
					m = f2(x)
				}
			case Reducer[R, T]: // reduce
				{
					count++
					if count == 1 {
						j = f2.initial
					}
					j = f2.reduce(j, x)
				}
			case Reducer[R, M]: // reduce
				{
					count++
					if count == 1 {
						j = f2.initial
					}
					j = f2.reduce(j, m)
				}
			case func() int: // take
				{
					num := f2()
					if i+1 >= num {
						ys = append(ys, m)
						goto exit
					}
				}
			default:
				return ys, j, fmt.Errorf("invalid function: %s", reflect.TypeOf(f))
			}
		}
		ys = append(ys, m)
	skip:
	}
exit:
	return ys, j, nil
}

type Reducer[K any, T any] struct {
	initial K
	reduce  func(acc K, item T) K
}

func Reduce[T any, K any](initial K, f func(acc K, item T) K) Reducer[K, T] {
	return Reducer[K, T]{
		initial: initial,
		reduce:  f,
	}
}
