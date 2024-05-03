package main

import (
	"fmt"
	"time"
)

type Task interface {
	Do()
}

type ActualTask struct {
	Done   chan struct{}
	Result int
	Err    error
}

func (at *ActualTask) Do() {
	defer func() {
		at.Done <- struct{}{}
	}()
	time.Sleep(time.Second)
	fmt.Println("task finish")
	panic("test")
	at.Result = 100
	at.Err = nil
}

type TaskPool struct {
	workernum int
	taskQueue chan Task
}

func NewTaskPool(workernum int) *TaskPool {
	return &TaskPool{
		workernum: workernum,
		taskQueue: make(chan Task, 10000),
	}
}

func (t *TaskPool) Start() {
	for range t.workernum {
		go t.startWorker()
	}
}

func (t *TaskPool) startWorker() {
	for task := range t.taskQueue {
		func() {
			defer func() {
				recover()
			}()
			task.Do()
		}()

	}
}

func (t *TaskPool) Stop() {
	close(t.taskQueue)
}

func (t *TaskPool) Commit(task Task) {
	t.taskQueue <- task
}

func main() {
	fmt.Println("hello task pool")

	tp := NewTaskPool(4)

	tp.Start()

	time.Sleep(time.Second)

	ts := []*ActualTask{}
	for range 10 {
		t := &ActualTask{Done: make(chan struct{}, 1)}
		tp.Commit(t)
		ts = append(ts, t)
	}

	for _, t := range ts {
		select {
		case <-t.Done:
			fmt.Println("done")
		}
	}

	for {
		time.Sleep(1 * time.Second)
	}
}
