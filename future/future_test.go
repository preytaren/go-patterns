package future

import (
	"testing"
	"time"
	"errors"
	"sync"
)

func TestBaseFuture_Get(t *testing.T) {
	f := NewBaseFuture(func(args ...interface{}) (interface{}, error) {
		if len(args) == 0 {
			return nil, errors.New("Invalid input")
		}
		time.Sleep(1 * time.Second)
		return args[0], nil
	}, func(args ...interface{}) string {
		return args[0].(string)
	})

	res, err := f.Get("hello", "world")
	if err != nil {
		t.Error(err)
	}

	if res != "hello" {
		t.Errorf("res not equal to hello, actual %s", res)
	}
}

func TestBaseFuture_Get2(t *testing.T) {
	f := NewBaseFuture(func(args ...interface{}) (interface{}, error) {
		if len(args) == 0 {
			return nil, errors.New("Invalid input")
		}
		time.Sleep(1 * time.Second)
		return args[0], nil
	}, func(args ...interface{}) string {
		return args[0].(string)
	})

	var wg sync.WaitGroup
	for i:=0; i<10; i++ {
		wg.Add(1)
		go func() {
			res, err := f.Get("hello", "world")
			if err != nil {
				t.Error(err)
			}

			if res != "hello" {
				t.Errorf("res not equal to hello, actual %s", res)
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func TestBaseFuture_GetError(t *testing.T) {
	f := NewBaseFuture(func(args ...interface{}) (interface{}, error) {
		if len(args) == 0 {
			return nil, errors.New("Invalid input")
		}
		time.Sleep(1 * time.Second)
		return args[0], errors.New("error")
	}, func(args ...interface{}) string {
		return args[0].(string)
	})

	var wg sync.WaitGroup
	for i:=0; i<10; i++ {
		wg.Add(1)
		go func() {
			_, err := f.Get("hello", "world")
			if err == nil {
				t.Error("Error should not be nil")
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func TestBaseFuture_Timeout(t *testing.T)  {
	f := NewBaseFuture(func(args ...interface{}) (interface{}, error) {
		if len(args) == 0 {
			return nil, errors.New("Invalid input")
		}
		time.Sleep(6 * time.Second)
		return args[0], errors.New("error")
	}, func(args ...interface{}) string {
		return args[0].(string)
	})
	var wg sync.WaitGroup
	for i:=0; i<10; i++ {
		wg.Add(1)
		go func() {
			_, err := f.Get("hello", "world")
			if err == nil {
				t.Error("Error should not be nil")
			}
			if err.Error() != "Get Timeout" {
				t.Error("Error is not a timeout", err)
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func BenchmarkBaseFuture_Get(b *testing.B) {
	f := NewBaseFuture(func(args ...interface{}) (interface{}, error) {
		if len(args) == 0 {
			return nil, errors.New("Invalid input")
		}
		for i:=0; i<10000000; i++ {}
		return args[0], nil
	}, func(args ...interface{}) string {
		return args[0].(string)
	})

	b.ReportAllocs()
	var wg sync.WaitGroup
	for i:=0; i<b.N; i++ {
		wg.Add(1)
		go func() {
			res, err := f.Get("hello", "world")
			if err != nil {
				b.Error(err)
			}

			if res != "hello" {
				b.Errorf("res not equal to hello, actual %s", res)
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func BenchmarkBase_Get2(b *testing.B) {
	f := func(args ...interface{}) (interface{}, error) {
		if len(args) == 0 {
			return nil, errors.New("Invalid input")
		}
		for i:=0; i<10000000; i++ {}
		return args[0], nil

	}
	var wg sync.WaitGroup
	for i:=0; i<b.N; i++ {
		wg.Add(1)
		go func() {
			f("hello", "world")
			wg.Done()
		}()
	}
	wg.Wait()
}
