package main

import (
    "fmt"
    "sync"
    "time"
)
func worker(id int, jobs <-chan int, wg *sync.WaitGroup) {
    defer wg.Done()
    for job := range jobs {
        fmt.Printf("Работяга %d делает работу %d\n", id, job)
        time.Sleep(time.Millisecond * 100)
    }
}
func main() {
    jobs := make(chan int, 10)
    var wg sync.WaitGroup
    for i := 1; i <= 3; i++ {
        wg.Add(1)
        go worker(i, jobs, &wg)
    }
    for i := 1; i <= 10; i++ {
        jobs <- i
    }
    close(jobs)
    wg.Wait()
    fmt.Println("Работа закончена!")
}
