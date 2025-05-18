package chan_pattern

import (
	"fmt"
	"sync"
	"time"

	"play_sync_lib/pkg"
)

func producer(ch chan<- int, items []int) {
	for _, item := range items {
		fmt.Printf("Producer: sending %d\n", item)
		ch <- item // チャネルに送信
	}
	close(ch) // 生産完了通知（全コンシューマーに伝わる）
	fmt.Println("Producer: closed the channel")
}

func consumer(id int, ch <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for item := range ch { // チャネルが閉じられるまで受信
		fmt.Printf("Consumer %d: received %d\n", id, item)
		time.Sleep(pkg.WaitTime * 2)
	}
	fmt.Printf("Consumer %d: channel closed, exiting\n", id)
}

func RunChanPattern() {
	ch := make(chan int, pkg.ConsumerNum) // バッファ付きチャネル
	var wg sync.WaitGroup

	// コンシューマー起動
	for i := 1; i <= pkg.ConsumerNum; i++ {
		wg.Add(1)
		go consumer(i, ch, &wg)
	}

	items := make([]int, pkg.ItemNum)
	for i := range items {
		items[i] = i + 1
	}

	// プロデューサー起動
	go producer(ch, items)

	wg.Wait()
}
