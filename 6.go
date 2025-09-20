package main
import (
    "fmt"
    "sync"
    "time"
)
type Logger struct {
    mu sync.Mutex
}
func (l *Logger) Log(message string) {
    l.mu.Lock()
    defer l.mu.Unlock()
    fmt.Printf("[%s] %s\n", time.Now().Format("15:04:05"), message)
}
func main() {
    logger := Logger{}
    var wg sync.WaitGroup
    for i := 0; i < 5; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            for j := 0; j < 3; j++ {
                logger.Log(fmt.Sprintf("Горутина %d: сообщение %d", id, j))
                time.Sleep(time.Millisecond * 10)
            }
        }(i)
    }
    wg.Wait()
}
