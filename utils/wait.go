package utils

import "sync"

type WaitGroup struct {
	Group *sync.WaitGroup
}

func InitWaitGroup() (*WaitGroup) {
	return &WaitGroup{
		Group:&sync.WaitGroup{},
	}
}

func (w *WaitGroup)Add(i int) {
	w.Group.Add(i)
	return
}

func (w *WaitGroup)Done() {
	w.Group.Done()
	return
}

func (w *WaitGroup)Wait() {
	w.Group.Wait()
	return
}
