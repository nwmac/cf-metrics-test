package main

import (
	"fmt"
	"io/ioutil"
	"runtime"
	"time"
)

func main() {

	fmt.Println("CF Metrics Test")
	fmt.Println("Running ...")

	delay := 24
	mem := 1
	doneFiles := 50

	blocks := make([][]byte, 32)

	for {

		// 8MB block
		block := make([]byte, 1024*1024*128)
		blocks[mem] = block
		done := make(chan int)

		// Write a file
		if doneFiles > 0 {
			file, err := ioutil.TempFile("", "cf-metrics-test")
			if err == nil {
				file.Truncate(1e7)
			} else {
				fmt.Println("Error creating file")
				fmt.Println(err)
			}
		}

		for i := 0; i < 2; i++ {
			go func(i int) {
				for {
					d := delay
					if d < 4 {
						d = 4
					}
					time.Sleep(time.Millisecond * time.Duration((d+10)/10))
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
			mem = mem + 1
		}

		fmt.Println(delay)
		if delay == -6 {
			delay = 24
			mem = 1
			blocks = make([][]byte, 32)
			fmt.Println("Resetting ...")
			runtime.GC()
		}

	}
}
