package ICashe

import (
	"testing"
	"time"
)

func TestICashe(t *testing.T) {
	var a ICache = InitTimeCash(10)
	a.AddWithTTL(-10, -10, 40*1000*20)
	for i := 0; i < 10; i++ {
		a.Add(i, i)
		time.Sleep(1000)
	}
	_, ok := a.Get(1)
	if !ok {
		t.Fatal("time")
	}
	for i := 10; i < 20; i++ {
		a.Add(i, i)
		time.Sleep(1000)
	}

	_, ok = a.Get(-10)
	if !ok {
		t.Fatal("time")
	}
	_, ok = a.Get(1)
	if ok {
		t.Fatal("time")
	}
	_, ok = a.Get(14)
	if !ok {
		t.Fatal("time")
	}
	time.Sleep(100000000)
	for i := 0; i < 20; i++ {
		a.Get(i)

	}
	a.Add(9, 9)
	_, ok = a.Get(-10)
	if ok {
		t.Fatal("time")
	}
	if !a.CheckIt() {
		t.Fatal("time")
	}
	n, ok := a.Get(9)
	if n != 9 || !ok {
		t.Fatal("time")
	}
	a.Add(9, 10)
	n, ok = a.Get(9)
	if n != 10 || !ok {
		t.Fatal("time")
	}

}
