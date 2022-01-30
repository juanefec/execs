package execs

import "sync"

// Executor takes an action and a handler that is executed after Do
// it returns Fire and Kill funtions to execute actions and handle them in a fully concurrent way
func Executor(do func(interface{}) interface{}, handle func(interface{})) (fire func(interface{}), kill func()) {
	var (
		in, out              = make(chan interface{}), make(chan interface{})
		killSend, killHandle = make(chan struct{}), make(chan struct{})
		wgSend, wgHandle     = sync.WaitGroup{}, sync.WaitGroup{}
	)

	go func() {
		for res := range out {
			wgHandle.Add(1)
			go func(i interface{}) {
				handle(i)
				wgHandle.Done()
			}(res)
		}
		killHandle <- struct{}{}
	}()

	go func() {
		for i := range in {
			wgSend.Add(1)
			go func(i interface{}) {
				out <- do(i)
				wgSend.Done()
			}(i)
		}
		killSend <- struct{}{}
	}()
	fire = func(i interface{}) { in <- i }
	kill = func() {
		close(in)
		<-killSend
		wgSend.Wait()
		close(out)
		<-killHandle
		wgHandle.Wait()
	}
	return
}
