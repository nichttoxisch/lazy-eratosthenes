package main

import "fmt"

func Lazy[V any](a V) func() V {
	return func() V { return a }
}

type Number func() int
type List[T any] []T
type Bool func() bool

func Sum(a, b Number) Number {
	return Lazy(a() + b())
}

func And(a, b Bool) Bool {
	if !a() {
		return Lazy(false)
	} else {
		return Lazy(b())
	}
}

func Or(a, b Bool) Bool {
	if a() {
		return Lazy(true)
	} else {
		return Lazy(b())
	}
}

func ToLazyList[T any](a List[T]) func() *LazyList[T] {
	return func() *LazyList[T] {
		if len(a) == 0 {
			return nil
		}

		return &LazyList[T]{
			Head: Lazy(a[0]),
			Tail: ToLazyList(a[1:]),
		}
	}
}

type LazyList[T any] struct {
	Head func() T
	Tail func() *LazyList[T]
}

func PrintList[T any](xs func() *LazyList[T]) {
	pair := xs()
	for pair != nil {
		fmt.Println(pair.Head())
		pair = pair.Tail()
	}
}

func Range(begin Number) func() *LazyList[int] {
	return func() *LazyList[int] {
		v := begin()
		return &LazyList[int]{
			Head: Lazy(v),
			Tail: Range(Lazy(v + 1)),
		}
	}
}

func Take[T any](n func() int, list func() *LazyList[T]) func() *LazyList[T] {
	return func() *LazyList[T] {
		m := n()
		if m == 0 {
			return nil
		}

		l := list()
		return &LazyList[T]{
			Head: l.Head,
			Tail: Take(Lazy(m-1), l.Tail),
		}
	}
}

func Filter(filter func(int) bool, list func() *LazyList[int]) func() *LazyList[int] {
	return func() *LazyList[int] {
		pair := list()
		if pair == nil {
			return nil
		}

		x := pair.Head()
		if filter(x) {
			return &LazyList[int]{
				Head: Lazy(x),
				Tail: Filter(filter, pair.Tail),
			}
		} else {
			return Filter(filter, pair.Tail)()
		}
	}
}

func Sieve(list func() *LazyList[int]) func() *LazyList[int] {
	return func() *LazyList[int] {
		pair := list()
		if pair == nil {
			return nil
		}
		y := pair.Head()
		return &LazyList[int]{
			Head: Lazy(y),
			Tail: Sieve(Filter(func(x int) bool { return x%y != 0 }, pair.Tail)),
		}
	}
}

func main() {
	prime := Sieve(Range(Lazy(2)))
	PrintList(prime)
}
