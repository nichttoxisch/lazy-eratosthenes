package main

import "fmt"

type LazyInt func() int64

type ListElement struct {
	head LazyInt
	tail LazyList
}
type LazyList func() ListElement

func count(begin LazyInt) LazyList {
	return func() ListElement {
		v := begin()
		return ListElement{
			head: func() int64 { return v },
			tail: count(func() int64 { return v + 1 }),
		}
	}
}

func printLazyList(xs LazyList) {
	for xs != nil {
		pair := xs()

		fmt.Println(pair.head())

		xs = pair.tail
	}
}

func take(n LazyInt, xs LazyList) LazyList {
	m := n()
	if m == 0 {
		return nil
	}

	return func() ListElement {
		pair := xs()
		return ListElement{
			head: pair.head,
			tail: take(func() int64 { return m - 1 }, pair.tail),
		}
	}
}

func filter(xs LazyList, f func(int64) bool) LazyList {
	if xs == nil {
		return nil
	}

	return func() ListElement {
		pair := xs()
		x := pair.head()
		if f(x) {
			return ListElement{
				head: pair.head,
				tail: filter(pair.tail, f),
			}
		} else {
			return filter(pair.tail, f)()
		}
	}
}

func sieve(xs LazyList) LazyList {
	if xs == nil {
		return nil
	}

	return func() ListElement {
		pair := xs()
		y := pair.head()
		return ListElement{
			head: func() int64 { return y },
			tail: sieve(filter(pair.tail, func(x int64) bool { return x%y != 0 })),
		}
	}
}

func main() {
	i2 := func() int64 { return 2 }

	prime := sieve(count(i2))

	printLazyList(
		prime,
	)
}
