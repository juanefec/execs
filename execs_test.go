package execs_test

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"testing"
	"time"

	"github.com/juanefec/execs"
)

func randomSleep() {
	nBig, err := rand.Int(rand.Reader, big.NewInt(200))
	if err != nil {
		panic(err)
	}
	n := nBig.Int64()
	time.Sleep(time.Millisecond * time.Duration(n))
}

type testClient struct {
	pipe  chan string
	store []string
	end   chan struct{}
}

func (tc *testClient) run() {
	for s := range tc.pipe {
		tc.store = append(tc.store, s)
	}
	tc.end <- struct{}{}
}

func (tc *testClient) Do(i interface{}) interface{} {
	msg, ok := i.(string)
	if !ok {
		return errors.New("haha u failed")
	}
	randomSleep()
	tc.pipe <- msg
	return "this was a succes!"
}

func TestT(t *testing.T) {
	tc := &testClient{make(chan string), make([]string, 0), make(chan struct{})}
	go tc.run()

	var (
		errPipe = make(chan error)
		errs    = make([]error, 0)
		okPipe  = make(chan string)
		oks     = make([]string, 0)
		end     = make(chan struct{})
	)

	go func() {
		for {
			select {
			case err := <-errPipe:
				errs = append(errs, err)
			case ok := <-okPipe:
				oks = append(oks, ok)
			case <-end:
				return
			}
		}
	}()

	handle := func(i interface{}) {
		switch r := i.(type) {
		case error:
			errPipe <- r
		case string:
			okPipe <- r
		}
	}

	fire, kill := execs.Executor(tc.Do, handle)
	wg := sync.WaitGroup{}
	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go func(i int) {
			fire(fmt.Sprintf("loool %v", i))
			wg.Done()
		}(i)
	}
	wg.Wait()
	kill()

	fire, kill = execs.Executor(tc.Do, handle)
	fire2, kill2 := execs.Executor(tc.Do, handle)
	for i := 0; i < 100000; i++ {
		fire(fmt.Sprintf("loool %v", i))
		fire2(fmt.Sprintf("loool %v", i))
	}

	for i := 0; i < 200; i++ {
		fire(10)
		fire2(10)
	}

	kill()
	kill2()

	close(tc.pipe)
	<-tc.end
	end <- struct{}{}

	equal(t, len(oks), len(tc.store))
	equal(t, len(tc.store), 210000)
	equal(t, len(oks), 210000)
	equal(t, len(errs), 400)
}

func equal(t *testing.T, e, v interface{}) {
	t.Helper()
	if e != v {
		t.Fail()
	}
}
