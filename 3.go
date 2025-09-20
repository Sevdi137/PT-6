package main

import (
    "fmt"
    "sync"
)

type TaskQ struct {
    tasks []string
    mu    sync.Mutex
}

func (q *TaskQ) Add(task string) {
    q.mu.Lock()
    defer q.mu.Unlock()
    q.tasks = append(q.tasks, task)
}

func (q *TaskQ) Get() string {
    q.mu.Lock()
    defer q.mu.Unlock()
    if len(q.tasks) == 0 {
        return ""
    }
    task := q.tasks[0]
    q.tasks = q.tasks[1:]
    return task
}

func main() {
    queue := TaskQ{}
    var wg sync.WaitGroup
    
    for i := 0; i < 3; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            for j := 0; j < 3; j++ {
                task := fmt.Sprintf("task-%d-%d", id, j)
                queue.Add(task)
            }
        }(i)
    }
    
    wg.Wait()
    
    for task := queue.Get(); task != ""; task = queue.Get() {
        fmt.Println("Processed:", task)
    }
}
