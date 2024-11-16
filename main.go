package main

import (
	"fmt"
	"log"
	"math"
	"sort"
	"sync"
	"time"
)

func isPrime(n int) bool {
	if n < 2 {
		return false
	}

	for i := 2; i <= int(math.Sqrt(float64(n))); i++ {
		if n%i == 0 {
			return false
		}
	}

	return true
}

func worker(tasks <-chan int, results chan<- int, progress *int, mu *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()
	for num := range tasks {
		if isPrime(num) {
			results <- num
		}
		mu.Lock()
		*progress++
		mu.Unlock()
	}
}

func execute() {
	const (
		numOfWorker = 4
		maxNum      = 100_000
	)

	var (
		wg       sync.WaitGroup
		tasks    = make(chan int, 1000)
		results  = make(chan int, 1000)
		progress int
		mu       = new(sync.Mutex)
	)

	for i := 1; i <= numOfWorker; i++ {
		wg.Add(1)
		go worker(tasks, results, &progress, mu, &wg)
	}

	go func() {
		ticker := time.NewTicker(1 * time.Millisecond)
		defer ticker.Stop()
		for range ticker.C {
			log.Printf("Processed %d / %d numbers ...\n", progress, maxNum)
		}
	}()

	go func() {
		for i := 1; i <= maxNum; i++ {
			tasks <- i
		}
		close(tasks)
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	primes := make([]int, 0)
	for prime := range results {
		primes = append(primes, prime)
	}

	sort.Ints(primes)

	fmt.Printf("Found %d primes.\n", len(primes))
	fmt.Println("Example primes:", primes[:10])
}

func main() {
	execute()
}
