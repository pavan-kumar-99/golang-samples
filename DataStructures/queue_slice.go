// Queue Data Structure
package main

import "fmt"

type queue []int

var prodqueue queue

func EnQueue(prodqueue queue, val int) queue {
	prodqueue = append(prodqueue, val)
	return prodqueue
}

func DeQueue(prodqueue queue) queue {
	return prodqueue[1:]
}

func main() {
	prodqueue = EnQueue(prodqueue, 1)
	prodqueue = EnQueue(prodqueue, 2)
	prodqueue = EnQueue(prodqueue, 3)
	fmt.Println("Before Dequee", prodqueue)
	prodqueue = DeQueue(prodqueue)
	fmt.Println("After Dequee", prodqueue)
}
// code scan please remove after testing
Access key ID: AKIAIOSFODNN7EXAMPLE

Secret access key: wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY
