# go-patterns (active)
一些go编码模式（持续更新）

## [future](https://github.com/preytaren/go-patterns/blob/master/future/future.go)
### 场景
在高并发场景下，可能存在大量重复的函数调用；无论是直接使用串行或并发执行，所有的调用都会被计算一次，导致调用的平均执行时间约等于单次调用时间；当单个调用很耗时时（如后端为db类存储，请求会被全部打到db），会严重影响系统的响应时间；实际上这些短时间内的重复调用只需要请求一次并把结果进行传递即可。
### 模式
对于一个调用周期内（从调用开始到结束），到达的其他它相同调用（获取资源相同）将被阻塞，直到第一个请求返回，并将这个值返回给这些阻塞请求。
对于一个调用周期内（从调用开始到结束），到达的其他它不同调用（获取资源不同）不会受到影响。
这个模式适用于单个调用耗时长（如IO或计算量大），且存在短期内重复调用，可以提高系统的平均响应时间和总吞吐量。
### 使用
首先创建一个future对象，参数分别为两个函数，第一个为实际future执行的函数f，第二个为从参数获取hashkey的函数；然后在future对象上正常进行函数调用即可。
执行future.Get(args) 等价于 f(args)，其中f为远函数
```
import "fmt"

func main() {
    f := NewBaseFuture(func(args ...interface{}) (interface{}, error) {
            if len(args) == 0 {
                return nil, errors.New("Invalid input")
            }
            time.Sleep(1 * time.Second)
            return args[0], nil
        }, func(args ...interface{}) string {
            return args[0].(string)
        })
    var wg = sync.WaitGroup
    for i:=0; i<10000; i++ {
         wg.Add(1)
         go func() {
             fmt.Printf("ARGS: %s\n" , f.Get("hello"))
             wg.Done()
         }()
    }
}
wg.Wait()
// 总计花费1s
```

## [pipeline](https://github.com/preytaren/go-patterns/blob/master/pipeline/pipeline.go)
### 场景

### 模式
使用pipeline模式，可以在pipeline端控制程序并发度；当底层实际执行有pipeline机制时，可以提升系统的吞吐量。
### 使用
使用Do添加pipeline执行函数及参数，调用pipeline.Sync执行整个pipeline，结果按照pipeline任务添加的顺序返回。
```
package main

func incr(a ...interface{}) (interface{}, error) {
	time.Sleep(1 * time.Microsecond)
	return a[0].(int) + 1, nil
}

func main() {
	bp := basicPipeliner{}
	for i := 0; i < 10; i++ {
		bp.Do(incr, i)
	}
	res, _ := bp.Sync()

}
```

## [config](https://github.com/preytaren/go-patterns/blob/master/config/config.go)
### 场景
使用config进行对象初始化。
### 模式
config模式用于通过config对象配置相应的参数, 使用With函数替代函数参数，可以拆分默认参数设置与对象初始化逻辑，提升代码可读性。
### 使用

```
cfg := NewConfig(1, WithTimeout(), WithC(10))
client := NewXX(cfg)
```

# TODO
## iterator
