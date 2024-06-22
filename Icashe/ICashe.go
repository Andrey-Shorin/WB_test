package ICashe

import (
	"container/heap"
	"sync"
	"time"
)

type ICache interface {
	CheckIt() bool
	Cap() int
	Clear()
	Add(key, value interface{})
	AddWithTTL(key, value interface{}, ttl time.Duration)
	Get(key interface{}) (value interface{}, ok bool)
	Remove(key interface{})
}

type Item struct {
	key interface{}
	t   time.Time
}

type PriorityQueue []Item

type TimeCash struct {
	mu       sync.Mutex
	m        map[interface{}]interface{}
	que      *PriorityQueue
	capacyty int
}

func InitTimeCash(Capacyty int) *TimeCash {
	tk := TimeCash{}
	tk.m = make(map[interface{}]interface{}, Capacyty)
	tk.que = new(PriorityQueue)
	heap.Init(tk.que)
	tk.capacyty = Capacyty
	return &tk
}
func (i *TimeCash) Cap() int {
	return i.que.Len()
}
func (i *TimeCash) CheckIt() bool {
	i.mu.Lock()
	defer i.mu.Unlock()
	if len(i.m) != i.que.Len() {
		return false
	}
	for _, j := range *i.que {
		_, ok := i.m[j.key]
		if !ok {
			return false
		}
	}
	return true
}
func (i *TimeCash) Clear() {
	i.mu.Lock()
	i.m = make(map[interface{}]interface{}, i.capacyty)
	i.que = new(PriorityQueue)
	heap.Init(i.que)
	i.mu.Unlock()
}

func (i *TimeCash) Add(key, value interface{}) {
	i.mu.Lock()
	i.add(key, value, time.Now())
	i.mu.Unlock()
}

func (i *TimeCash) AddWithTTL(key, value interface{}, ttl time.Duration) {
	i.mu.Lock()
	i.add(key, value, time.Now().Add(ttl))
	i.mu.Unlock()
}

func (i *TimeCash) add(key, value interface{}, t time.Time) {
	_, ok := i.m[key]
	if ok {
		i.m[key] = value
		j := i.que.find(key)
		(*(i.que))[j].t = t
		return
	}
	if i.que.Len() < i.capacyty {
		i.m[key] = value
		heap.Push(i.que, Item{key, t})
		return
	}
	old := (heap.Pop(i.que)).(Item)
	if old.t.After(t) {
		heap.Push(i.que, old)
		return
	}
	//fmt.Print("remove  -- ")
	//fmt.Println(h)
	delete(i.m, old.key)
	i.m[key] = value
	heap.Push(i.que, Item{key, t})
}
func (i *TimeCash) Get(key interface{}) (value interface{}, ok bool) {
	i.mu.Lock()
	defer i.mu.Unlock()
	value, ok = i.m[key]
	if ok {
		j := i.que.find(key)
		t := (*(i.que))[j].t
		if t.After(time.Now()) {
			return value, ok
		}
		(*(i.que))[j].t = time.Now()
		heap.Fix(i.que, j)
		return value, ok
	}
	return value, ok

}
func (i *TimeCash) Remove(key interface{}) {
	i.mu.Lock()
	_, ok := i.m[key]
	if ok {
		j := i.que.find(key)
		heap.Remove(i.que, j)
		delete(i.m, key)
	}
	i.mu.Unlock()

}
func (q *PriorityQueue) find(key interface{}) int {
	for i, val := range *q {
		if val.key == key {
			return i
		}
	}
	//panic("n not equal") // may be i.Clear() ???
	return -1
}

func (h *PriorityQueue) Push(x any) {
	*h = append(*h, x.(Item))
}

func (h PriorityQueue) Len() int           { return len(h) }
func (h PriorityQueue) Less(i, j int) bool { return h[i].t.Before(h[j].t) }
func (h PriorityQueue) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *PriorityQueue) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
