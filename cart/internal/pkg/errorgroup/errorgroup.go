// Package errorgroup Упрощённый аналог errgroup
package errorgroup

import (
	"sync"
)

// Group Группа
type Group struct {
	wg       *sync.WaitGroup
	errMutex *sync.RWMutex
	err      error
}

// New Конструктор
func New() *Group {
	return &Group{
		wg:       &sync.WaitGroup{},
		errMutex: &sync.RWMutex{},
	}
}

// Go Запуск горутины
func (g *Group) Go(f func() error) {
	g.errMutex.RLock()
	defer g.errMutex.RUnlock()
	if g.err != nil {
		return
	}

	g.wg.Add(1)
	go func() {
		defer g.wg.Done()
		if err := f(); err != nil {
			g.errMutex.Lock()
			defer g.errMutex.Unlock()
			if g.err == nil {
				g.err = err
			}
		}
	}()
}

// Wait Ожидание выполнения
func (g *Group) Wait() error {
	g.wg.Wait()
	return g.err
}
