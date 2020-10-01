package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

func main() {

	fmt.Println("CF Metrics Test")
	fmt.Println("Running ...")
	reportDiskUsage()

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
			file, err := ioutil.TempFile(".", "cf-metrics-test")
			if err == nil {
				// fmt.Println(file.Name())
				file.Truncate(1024 * 1024 * 10)
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

		// Report disk usage
		reportDiskUsage()

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

func reportDiskUsage() {
	size, err := dirSize(".")
	if err == nil {
		fmt.Printf("Disk usage: %s\n", byteCountSI(size))
	}
}

func dirSize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})
	return size, err
}

func byteCountSI(b int64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB",
		float64(b)/float64(div), "kMGTPE"[exp])
}
