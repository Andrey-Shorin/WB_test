package main

import (
	"fmt"
	"time"

	ICashe "main/Icashe"
)

func main() {
	var a ICashe.ICache = ICashe.InitTimeCash(10)
	fmt.Printf("cap = %d\n", a.Cap())
	a.AddWithTTL(-10, -10, 40*1000*20)
	for i := 0; i < 10; i++ {
		a.Add(i, i)
		time.Sleep(1000)
	}
	fmt.Println(a.Get(1))
	for i := 10; i < 20; i++ {
		a.Add(i, i)
		time.Sleep(1000)
	}
	fmt.Printf("cap = %d\n", a.Cap())
	for i := 0; i < 20; i++ {
		fmt.Println(a.Get(i))
		time.Sleep(1000)
	}
	fmt.Printf("cap = %d\n", a.Cap())
	fmt.Println(a.Get(-10))
	fmt.Println(a.Cap())

	time.Sleep(100000000)
	a.Add(9, 9)
	fmt.Println(a.Get(-10))
	fmt.Println(a.Get(9))
}
