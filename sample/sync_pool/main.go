package main

import (
	"bytes"
	"fmt"
	"sync"
	"time"
)

type Client struct {
	buf bytes.Buffer
	seq int
}

func (c *Client) Do(msg string) {
	c.seq++

	c.buf.Reset()
	c.buf.WriteString(fmt.Sprintf("[%d] %s", c.seq, msg))

	// 模擬 IO
	time.Sleep(10 * time.Millisecond)

	fmt.Println(c.buf.String())
}

var pool = sync.Pool{
	New: func() any {
		fmt.Println("Create Client")
		return &Client{}
	},
}

func main() {

	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {

		wg.Add(1)

		go func(id int) {
			defer wg.Done()

			client := pool.Get().(*Client)
			defer pool.Put(client)

			client.Do(fmt.Sprintf("request-%d", id))

		}(i)
	}

	wg.Wait()
}