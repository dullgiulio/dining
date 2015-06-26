package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

var take chan struct{}

func init() {
	rand.Seed(time.Now().UnixNano())
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func delay(maxRand int) {
	time.Sleep(50*time.Millisecond + time.Duration(rand.Intn(maxRand)))
}

type fork struct {
	sync.Mutex
	id int
}

func (f *fork) String() string {
	return fmt.Sprintf("fork %d", f.id)
}

func newForks(n int) []*fork {
	c := make([]*fork, n)
	for i := 0; i < n; i++ {
		c[i] = &fork{id: i}
	}
	return c
}

type philosopher struct {
	name        string
	left, right *fork
}

func newPhilosopher(fs []*fork, name string, left, right int) *philosopher {
	return &philosopher{name: name, left: fs[left], right: fs[right]}
}

func (p *philosopher) think() {
	delay(10)
}

func (p *philosopher) eat() {
	take <- struct{}{}
	p.left.Lock()
	p.right.Lock()

	fmt.Printf("%s started eating\n", p.name)
	time.Sleep(100 * time.Millisecond)
	fmt.Printf("%s finished eating\n", p.name)

	p.left.Unlock()
	p.right.Unlock()
	<-take
}

func main() {
	fs := newForks(4)
	ps := []*philosopher{
		newPhilosopher(fs, "Baruch Spinoza", 0, 1),
		newPhilosopher(fs, "Gilles Deleuze", 1, 2),
		newPhilosopher(fs, "Karl Marx", 2, 3),
		newPhilosopher(fs, "Friedrich Nietzsche", 3, 0),
	}
	take = make(chan struct{}, len(ps)-1)
	wg := &sync.WaitGroup{}

	wg.Add(len(ps))
	for _, p := range ps {
		go func(p *philosopher) {
			p.think()
			p.eat()
			wg.Done()
		}(p)
	}

	wg.Wait()
}
