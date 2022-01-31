package execs

import "sync"

// Executor takes an action and a handler that is executed after the action
// it returns Fire and Kill funtions to execute the action and handle them in a fully concurrent way
func Executor(action func(interface{}) interface{}, handle func(interface{})) (fire func(interface{}), kill func()) {
	var (
		in, out              = make(chan interface{}), make(chan interface{})
		killSend, killHandle = make(chan struct{}, 1), make(chan struct{}, 1)
		wgSend, wgHandle     = sync.WaitGroup{}, sync.WaitGroup{}

		sender = func() {
			for i := range in {
				wgSend.Add(1)
				go func(i interface{}) {
					out <- action(i)
					wgSend.Done()
				}(i)
			}
			killSend <- struct{}{}
		}

		handler = func() {
			for res := range out {
				wgHandle.Add(1)
				go func(i interface{}) {
					handle(i)
					wgHandle.Done()
				}(res)
			}
			killHandle <- struct{}{}
		}
	)

	fire = func(i interface{}) { in <- i }

	kill = func() {
		close(in)
		<-killSend
		wgSend.Wait()
		close(out)
		<-killHandle
		wgHandle.Wait()
	}

	go handler()

	go sender()

	return
}
