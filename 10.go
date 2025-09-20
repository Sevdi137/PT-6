package main
import (
    "errors"
    "fmt"
    "sync"
)
type ErrColl struct {
    errors []error
    mu     sync.Mutex
}

func (ec *ErrColl) Add(err error) {
    ec.mu.Lock()
    defer ec.mu.Unlock()
    ec.errors = append(ec.errors, err)
}

func (ec *ErrColl) GetErrors() []error {
    ec.mu.Lock()
    defer ec.mu.Unlock()
    return ec.errors
}

func stage1(input <-chan int, output chan<- int, ec *ErrColl, wg *sync.WaitGroup) {
    defer wg.Done()
    for num := range input {
        if num%5 == 0 {
            ec.Add(errors.New("Стадия 1: число делимое на 5"))
            continue
        }
        output <- num * 2
    }
    close(output)
}

func stage2(input <-chan int, output chan<- int, ec *ErrColl, wg *sync.WaitGroup) {
    defer wg.Done()
    for num := range input {
        if num%7 == 0 {
            ec.Add(errors.New("Стадия 2: число делимое на 7"))
            continue
        }
        output <- num + 1
    }
    close(output)
}

func stage3(input <-chan int, output chan<- int, ec *ErrColl, wg *sync.WaitGroup) {
    defer wg.Done()
    for num := range input {
        result := num - 3
        if result < 0 {
            ec.Add(fmt.Errorf("Стадия 3: негативный результат %d", result))
            continue
        }
        output <- result
    }
    close(output)
}

func main() {
    input := make(chan int, 10)
    stage1out := make(chan int, 10)
    stage2out := make(chan int, 10)
    results := make(chan int, 10)
    var ec ErrColl
    var wg sync.WaitGroup
    wg.Add(3)
    go stage1(input, stage1out, &ec, &wg)
    go stage2(stage1out, stage2out, &ec, &wg)
    go stage3(stage2out, results, &ec, &wg)
    for i := 1; i <= 20; i++ {
        input <- i
    }
    close(input)

    go func() {
        for result := range results {
            fmt.Printf("Результат: %d\n", result)
        }
    }()
    wg.Wait()
		fmt.Println("Ошибки:")
    for _, err := range ec.GetErrors() {
        fmt.Println(err)
    }
}
