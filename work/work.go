package work

import (
	"sync"
)

// work包的目的是展示如何使用无缓冲的通道来创建一个goroutine池, 这些goroutine执行并控制一组工作
// 让其并发工作。在这种情况下, 使用无缓冲的通道要比随意指定一个缓冲区大小的有缓冲的通道好，
// 因为这种情况既不需要一个工作队列，也不需要一组goroutine配合执行。
// 无缓冲的通道保证两个goroutine之间的数据交换。这种使用无缓冲的通道的方法允许调用者知道
// 什么时候goroutine池正在执行工作, 而且如果池中的所有goroutine都在忙，也可以及时通过通道通知调用者。
// 使用无缓冲的通道不会有工作在队列中丢失或者卡住，所有工作都会被处理。

// Worker必须满足接口类型
// 才能使用工作池
type Worker interface {
	Task()
}

// Pool提供一个goroutine池, 这个池可以完成任何已提交的Worker任务
type Pool struct {
	work chan Worker
	wg   sync.WaitGroup
}

// New创建一个新协程池
func New(maxGoroutines int) *Pool {
	p := Pool{
		work: make(chan Worker),
	}
	p.wg.Add(maxGoroutines)
	for i := 0; i < maxGoroutines; i++ {
		go func() {
			// p.work关闭的时候将该协程从waitgroup中关闭
			defer p.wg.Done()
			for w := range p.work {
				// 阻塞等待执行任务
				w.Task()
			}
		}()
	}
	return &p
}

// Run提交工作到协程池
func (p *Pool) Run(w Worker) {
	p.work <- w
}

// Shutdown 等待所有goroutine停止工作
func (p *Pool) Shutdown() {
	close(p.work)
	p.wg.Wait()
}
