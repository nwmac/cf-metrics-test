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

	delay := 32
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

		delay = delay - 1
		if delay >= 0 {
			cpus = cpus + 1
			mem = mem + 1
		}

		fmt.Println(delay)
		if delay == -5 {
			delay = 32
			cpus = 1
			mem = 1
			blocks = make(map[int]interface{})
			fmt.Println("Resetting ...")
		}

	}

	fmt.Println("CF Metrics Test finished")

}
