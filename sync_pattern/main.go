package sync_pattern

import (
	"fmt"
	"sync"
	"time"

	"play_sync_lib/pkg"
)

type QueueManager struct {
	mu    sync.Mutex
	cond  *sync.Cond
	queue []int
	done  bool
}

func NewQueueManager() *QueueManager {
	qm := &QueueManager{}
	qm.cond = sync.NewCond(&qm.mu)
	return qm
}

func (q *QueueManager) Produce(items []int) {
	for _, item := range items {
		q.mu.Lock()
		fmt.Printf("Producer: Adding %d to queue...\n", item)
		q.queue = append(q.queue, item)
		q.cond.Broadcast() // 通知は1つだけ（複数なら Broadcast）
		q.mu.Unlock()
	}

	q.mu.Lock()
	q.done = true
	q.cond.Broadcast() // 全コンシューマに終了を通知
	q.mu.Unlock()
}

func (q *QueueManager) Consume(id int) {
	for {
		q.mu.Lock()
		for len(q.queue) == 0 && !q.done {
			fmt.Printf("Consumer %d: Waiting for data...\n", id)
			q.cond.Wait()
		}

		if len(q.queue) == 0 && q.done {
			q.mu.Unlock()
			fmt.Printf("Consumer %d: No more data. Exiting.\n", id)
			return
		}

		item := q.queue[0]
		q.queue = q.queue[1:]
		time.Sleep(pkg.WaitTime * 2)
		fmt.Printf("Consumer %d: Consumed %d, Queue: %v\n", id, item, q.queue)
		q.mu.Unlock()
	}
}

func RunSyncPattern() {
	qm := NewQueueManager()
	var wg sync.WaitGroup

	// コンシューマー起動
	for i := 1; i <= pkg.ConsumerNum; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			qm.Consume(id)
		}(i)
	}

	items := make([]int, pkg.ItemNum)
	for i := range items {
		items[i] = i + 1
	}

	// プロデューサー起動
	wg.Add(1)
	go func() {
		defer wg.Done()
		qm.Produce(items)
	}()

	wg.Wait()
}
