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
			//case func(T) K:
			//	{
			//		x = f2(x)
			//	}
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
					//fmt.Printf("-----\n")
					//fmt.Printf("ORIGINAL 1: '%v'\n", x)
					x = f2(x)
					//fmt.Printf("MAPPED 1: '%v'\n", x)
					//fmt.Printf("-----\n")
					if v, ok := any(x).(M); ok {
						m = v
					}
				}
			case func(T) M: // map
				{
					//fmt.Printf("-----\n")
					//fmt.Printf("ORIGINAL 2: '%v'\n", x)
					m = f2(x)
					//fmt.Printf("MAPPED 2: '%v'\n", m)
					//fmt.Printf("-----\n")
				}
			case reduceStruct[R, T]: // reduce
				{
					count++
					//fmt.Printf("REDUCE (original): '%v'\n", x)
					if count == 1 {
						j = f2.initial
						//f2.initial = j
					}
					j = f2.reduce(j, x)
					//j = f2.initial
				}
			case reduceStruct[R, M]: // reduce
				{
					count++
					//fmt.Printf("REDUCE (mapped): '%v'\n", m)
					//f2.initial = j
					if count == 1 {
						j = f2.initial
						//f2.initial = j
					}
					j = f2.reduce(j, m)
					//j = f2.initial
				}
			case func() int: // take
				{
					num := f2()
					if i+1 >= num {
						ys = append(ys, m)
						goto exit
					}
					//f2.reduce(f2.initial, x)
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

type reduceStruct[K any, T any] struct {
	initial K
	reduce  func(acc K, item T) K
}

func Reduce[T any, K any](initial K, f func(acc K, item T) K) reduceStruct[K, T] {
	return reduceStruct[K, T]{
		initial: initial,
		reduce:  f,
	}
}
