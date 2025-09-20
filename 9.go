package main
import (
    "fmt"
    "sync"
    "time"
)
type Metrics struct {
    SucReq int
    FailReq int
    ResTime time.Duration
    mu sync.Mutex
}
func (m *Metrics) IncrementSuccess() {
    m.mu.Lock()
    defer m.mu.Unlock()
    m.SucReq++
}
func (m *Metrics) IncrementFailure() {
    m.mu.Lock()
    defer m.mu.Unlock()
    m.FailReq++
}
func (m *Metrics) AddResTime(d time.Duration) {
    m.mu.Lock()
    defer m.mu.Unlock()
    m.ResTime += d
}
func (m *Metrics) Report() {
    m.mu.Lock()
    defer m.mu.Unlock()
    fmt.Println("Отчет метрики")
    fmt.Printf("Успешный запрос: %d\n", m.SucReq)
    fmt.Printf("Неудачный запрос: %d\n", m.FailReq)
    fmt.Printf("Общее время: %v\n", m.ResTime)
    if m.SucReq > 0 {
        fmt.Printf("Среднее время: %v\n", m.ResTime/time.Duration(m.SucReq))
    }
}
func main() {
    var metrics Metrics
    var wg sync.WaitGroup
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            for j := 0; j < 5; j++ {
                if j%3 == 0 {
                    metrics.IncrementFailure()
                } else {
                    metrics.IncrementSuccess()
                    metrics.AddResTime(time.Duration(100+j*10) * time.Millisecond)
                }
                time.Sleep(time.Millisecond * 50)
            }
        }(i)
    }
    wg.Wait()
    metrics.Report()
}
