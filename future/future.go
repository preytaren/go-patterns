/*
parrallel Get with same params on future will only be execute once, the rest will be blocked until the first requests return.
Future can speed up parrallel execution
 */
package future

import (
	"time"
	"sync"
	"errors"
)

const (
	FUTURETIMEOUT = 5 * time.Second
)

type Future interface {
	Get(...interface{}) (interface{}, error)
}

type BaseFuture struct {
	sync.RWMutex
	inner  map[string][]chan msg
	data   map[string]interface{}

	// caller function, which should be time consumption
	f      func(...interface{}) (interface{}, error)

	// function that convert args to a string
	keyF   func(...interface{}) string
	ticker *time.Ticker
}

type msg struct {
	raw interface{}
	err error
}

func NewBaseFuture(f func(...interface{}) (interface{}, error), keyF func(...interface{}) string) Future {
	sf := new(BaseFuture)
	sf.RWMutex = sync.RWMutex{}
	sf.inner = make(map[string][]chan msg)
	sf.f = f
	sf.keyF = keyF
	return sf
}

func (sf *BaseFuture) Get(args ...interface{}) (interface{}, error) {
	key := sf.keyF(args...)
	sf.Lock()
	_, ok := sf.inner[key]
	if ok {
		ch := make(chan msg, 1)
		sf.inner[key] = append(sf.inner[key], ch)
		sf.Unlock()
		ticker := time.NewTicker(FUTURETIMEOUT)
		defer ticker.Stop()
		select {
		case r :=<- ch:
			return r.raw, r.err
		case <- ticker.C:
			return nil, errors.New("Get Timeout")
		}
	} else {
		sf.inner[key] = make([]chan msg, 0, 4)
		sf.Unlock()
		ch := make(chan msg, 1)
		tick := time.NewTicker(FUTURETIMEOUT)
		defer tick.Stop()
		go func() {
			res, err := sf.f(key)
			ch <- msg{res, err}
		}()
		var res msg
		select {
		case res = <- ch:
		case <-tick.C:
			res = msg{nil, errors.New("Get Timeout")}
		}
		sf.Lock()
		chList := sf.inner[key]
		delete(sf.inner, key)
		sf.Unlock()
		if res.err != nil {
			go func() {
				for _, ch := range chList {
					ch <- res
					close(ch)
				}
			}()
			return nil, res.err
		}
		go func() {
			for _, ch := range chList {
				ch <- res
				close(ch)
			}
		}()
		return res.raw, nil
	}
}

