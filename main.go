package main

import (
	"fmt"
	"time"
)

func main() {

	fmt.Println("CF Metrics Test")
	fmt.Println("Starting in 10 seconds ...")

	//time.Sleep(time.Second * 10)
	fmt.Println("Running ...")

	delay := 20
	cpus := 1
	mem := 1

	blocks := make(map[int]interface{})

	for {

		// 8MB block
		block := make([]byte, mem*1024*1024*64)
		blocks[mem] = block
		done := make(chan int)

		for i := 0; i < cpus; i++ {
			go func(i int) {
				for {
					select {
					case <-done:
						return
					default:
					}
				}
			}(i)
		}

		time.Sleep(time.Second * 10)
		close(done)

		cpus = cpus + 1
		mem = mem + 1

		if mem > 32 {
			mem = 1
			blocks = make(map[int]interface{})

		}

		delay = delay - 1
		fmt.Println(delay)
		if delay == -1 {
			break
		}

	}

	fmt.Println("CF Metrics Test finished")

}
