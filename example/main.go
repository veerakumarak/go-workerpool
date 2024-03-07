package main

import (
	"fmt"
	"github.com/veerakumarak/go-workerpool"
	"time"
)

func CreateTask(i int) func() {
	return func() {
		fmt.Println("task started ", i)
		time.Sleep(time.Second * 5)
		fmt.Println("task completed ", i)
	}
}

func main() {
	pool := workerpool.New("default", 10, 1000)
	pool.Start()

	for i := 0; i < 4; i++ {
		newTask := CreateTask(i)
		pool.Submit(newTask)
	}

	pool.Shutdown()
}
