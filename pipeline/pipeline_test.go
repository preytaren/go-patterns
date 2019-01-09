package pipeline

import (
	"sync"
	"testing"
	"time"
)

const TotalCount = 305

func incr(a ...interface{}) (interface{}, error) {
	time.Sleep(1 * time.Microsecond)
	return a[0].(int) + 1, nil
}

func TestBasicPipeliner_Do(t *testing.T) {
	bp := basicPipeliner{}
	for i := 0; i < 10; i++ {
		bp.Do(incr, i)
	}
	res, _ := bp.Sync()
	for i, r := range res {
		if r.(int) != i+1 {
			t.Errorf("pipeline %d should equal to %d", r, i)
		}
	}
}

func BenchmarkBasicPipeliner_Do(b *testing.B) {
	bp := basicPipeliner{}
	b.ReportAllocs()
	for j := 0; j < b.N; j++ {
		for i := 0; i < TotalCount; i++ {
			bp.Do(incr, i)
		}
		bp.Sync()
	}
}

func TestBondPipeliner_Do(t *testing.T) {
	bp := NewBondPipeliner(6)
	for i := 0; i < 10; i++ {
		bp.Do(incr, i)
	}
	res, _ := bp.Sync()
	for i, r := range res {
		if r.(int) != i+1 {
			t.Errorf("pipeline %d should equal to %d", r, i)
		}
	}
}

func TestBondPipeliner_Do2(t *testing.T) {
	bp := NewBondPipeliner(12)
	for i := 0; i < 10; i++ {
		bp.Do(incr, i)
	}
	res, _ := bp.Sync()
	for i, r := range res {
		if r.(int) != i+1 {
			t.Errorf("pipeline %d should equal to %d", r, i)
		}
	}
}

func BenchmarkBondPipeliner_Do(b *testing.B) {
	bp := NewBondPipeliner(TotalCount - 1)
	b.ReportAllocs()
	for j := 0; j < b.N; j++ {
		for i := 0; i < TotalCount; i++ {
			bp.Do(incr, i)
		}
		bp.Sync()
	}
}

func BenchmarkNoPipeliner(b *testing.B) {
	b.ReportAllocs()
	var wg sync.WaitGroup
	for j := 0; j < b.N; j++ {
		for i := 0; i < TotalCount; i++ {
			wg.Add(1)
			go func() {
				incr(j)
				wg.Done()
			}()
		}
	}
	wg.Wait()
}
