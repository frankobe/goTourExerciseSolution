package main

import (
	"fmt"
	"golang.org/x/tour/tree"
)

// Make sure go get golang.org/x/tour/tree first if running on local

// Walk walks the tree t sending all values
// from the tree to the channel ch.
func Walk(t *tree.Tree, ch chan int) {
	walkHelper(t, ch)
	close(ch)
}

func walkHelper(t *tree.Tree, ch chan int) {
	if t != nil {
		walkHelper(t.Left, ch)
		ch <- t.Value
		walkHelper(t.Right, ch)
	}
}

// Same determines whether the trees
// t1 and t2 contain the same values.
func Same(t1, t2 *tree.Tree) bool {
	ch1 := make(chan int)
	ch2 := make(chan int)

	go Walk(t1, ch1)
	go Walk(t2, ch2)
	for {
		i, cont1 := <-ch1
		j, cont2 := <-ch2

		if (i != j) || (cont1 != cont2) {
			return false
		}

		if cont1 == cont2 == true {
			break
		}
	}
	return true
}

func main() {
	//  ch := make(chan int)
	//  go Walk(tree.New(2), ch)
	//  for v := range ch {
	//    fmt.Print(v)
	//  }
	fmt.Println(Same(tree.New(1), tree.New(1)))
	fmt.Println(Same(tree.New(1), tree.New(2)))
}
