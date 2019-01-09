package pipeline

import (
	"sync"
	"sync/atomic"
)

type Pipeliner interface {
	Do(func(...interface{}) (interface{}, error), ...interface{})
	Sync() ([]interface{}, []error)
}

type block struct {
	f    func(...interface{}) (interface{}, error)
	args []interface{}
}

type result struct {
	raw interface{}
	err error
}

type basicPipeliner struct {
	blocks []*block
	synced int32
	sync.Mutex
}

func (bp *basicPipeliner) Do(f func(...interface{}) (interface{}, error), args ...interface{}) {
	if atomic.LoadInt32(&bp.synced) == 0 {
		bp.blocks = append(bp.blocks, &block{f, args})
	}
}

func (bp *basicPipeliner) Sync() ([]interface{}, []error) {
	if atomic.CompareAndSwapInt32(&bp.synced, 0, 1) {
		ret, errs := pipeSync(bp.blocks)
		bp.blocks = bp.blocks[:0]
		atomic.StoreInt32(&bp.synced, 0)
		return ret, errs
	} else {
		return nil, nil
	}
}

func pipeSync(blocks []*block) ([]interface{}, []error) {
	rs := make([]*result, len(blocks))
	var wg sync.WaitGroup
	for i, b := range blocks {
		i := i
		b := b
		wg.Add(1)
		go func() {
			res, err := b.f(b.args...)
			rs[i] = &result{res, err}
			wg.Done()
		}()
	}
	wg.Wait()
	ret := make([]interface{}, len(rs))
	errs := make([]error, len(rs))
	for i, _ := range rs {
		ret[i] = rs[i].raw
		errs[i] = rs[i].err
	}
	return ret, errs
}

type ChuckedPipeliner struct {
	blocks    []*block
	synced    int32
	bondCount int
}

func NewBondPipeliner(sizeBond int) Pipeliner {
	return &ChuckedPipeliner{bondCount: sizeBond}
}

func (bp *ChuckedPipeliner) Do(f func(...interface{}) (interface{}, error), args ...interface{}) {
	if atomic.LoadInt32(&bp.synced) == 0 {
		bp.blocks = append(bp.blocks, &block{f, args})
	}
}

func (bp *ChuckedPipeliner) Sync() ([]interface{}, []error) {
	if atomic.CompareAndSwapInt32(&bp.synced, 0, 1) {
		var ret []interface{}
		var errs []error
		if len(bp.blocks) <= bp.bondCount {
			ret, errs := pipeSync(bp.blocks)
			bp.blocks = bp.blocks[:0]
			atomic.StoreInt32(&bp.synced, 0)
			return ret, errs
		} else {
			offset := 0
			for offset < len(bp.blocks) {
				var res []interface{}
				var err []error
				if offset+bp.bondCount >= len(bp.blocks) {
					res, err = pipeSync(bp.blocks[offset:])
				} else {
					res, err = pipeSync(bp.blocks[offset : offset+bp.bondCount])
				}
				ret = append(ret, res...)
				errs = append(errs, err...)
				offset += bp.bondCount
			}
			return ret, errs
		}
	}
	return nil, nil
}
