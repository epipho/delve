package main

import "fmt"

type FooBar struct {
	Baz int
	Bur string
}

// same member names, different order / types
type FooBar2 struct {
	Bur int
	Baz string
}

func barfoo() {
	a1 := "bur"
	fmt.Println(a1)
}

func foobar(baz string, bar FooBar) {
	var (
		a1  = "foofoofoofoofoofoo"
		a2  = 6
		a3  = 7.23
		a4  = [2]int{1, 2}
		a5  = []int{1, 2, 3, 4, 5}
		a6  = FooBar{Baz: 8, Bur: "word"}
		a7  = &FooBar{Baz: 5, Bur: "strum"}
		a8  = FooBar2{Bur: 10, Baz: "feh"}
		a9  = (*FooBar)(nil)
		a10 = a1[2:5]
		neg = -1
		i8  = int8(1)
		f32 = float32(1.2)
		i32 = [2]int32{1, 2}
	)

	barfoo()
	fmt.Println(a1, a2, a3, a4, a5, a6, a7, a8, a9, a10, baz, neg, i8, f32, i32, bar)
}

func main() {
	foobar("bazburzum", FooBar{Baz: 10, Bur: "lorem"})
}
